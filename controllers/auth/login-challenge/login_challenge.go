package main

import (
	"password-caddy/api/core/container"
	coreTypes "password-caddy/api/core/types"
	"password-caddy/api/lib/dynamoclient"
	"password-caddy/api/lib/logger"
	"password-caddy/api/lib/result"
	"password-caddy/api/lib/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sesTypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type LoginChallengeRequest struct {
	Email string
	Code  string
}

// Initialize the Login Challenge Request
func Init(event events.APIGatewayProxyRequest) *result.Result {
	code, _ := util.GenerateOTP(6)
	email := event.PathParameters["email"]

	return result.SuccessWithValue(
		200,
		LoginChallengeRequest{
			Email: email,
			Code:  code,
		},
	)
}

// Get the SES verification status of the email address
func GetEmailStatus(res result.ResultValue) *result.Result {
	request := res.(LoginChallengeRequest)

	response := container.SesClient().
		GetVerificationStatus(request.Email)

	if !response.IsSuccess {
		logger.Error(
			"Failed to get the verification status of email address",
			struct {
				Email string
				Error coreTypes.PasswordCaddyError
			}{
				Email: request.Email,
				Error: response.Error,
			},
		)

		return result.Failure(
			response.Error.StatusCode,
			response.Error.Message,
		)
	}

	data := response.Data.(struct{ Status sesTypes.VerificationStatus })

	if data.Status != "Success" {
		logger.Warn(
			"Attempted to challenge an email that is not verified",
			struct {
				Email              string
				VerificationStatus sesTypes.VerificationStatus
			}{
				Email:              request.Email,
				VerificationStatus: data.Status,
			},
		)

		return result.Failure(401, "Unauthorized to challenge email address")
	}

	logger.Info(
		"Successfully verified that email address can be challenged",
		struct {
			Email              string
			VerificationStatus sesTypes.VerificationStatus
		}{
			Email:              request.Email,
			VerificationStatus: data.Status,
		},
	)

	return result.SuccessWithValue(200, request)
}

// Update the user status in DynamoDB to ACTIVE
func UpdateEmailStatusInDynamo(res result.ResultValue) *result.Result {
	request := res.(LoginChallengeRequest)

	dynamoRequest := dynamoclient.DyanamoUpdateRequest{
		Key: request.Email,
		Values: map[string]dynamoclient.DynamoUpdateItem{
			"STATUS": {
				Action: types.AttributeActionPut,
				Value:  "ACTIVE",
			},
		},
	}

	response := container.DynamoClient().
		Update(dynamoRequest)

	if !response.IsSuccess {
		logger.Error(
			"Failed to update email status in DynamoDB",
			struct {
				Email string
				Error coreTypes.PasswordCaddyError
			}{
				Email: request.Email,
				Error: response.Error,
			},
		)

		return result.Failure(
			response.Error.StatusCode,
			response.Error.Message,
		)
	}

	logger.Info(
		"Updated email status in DyanamoDB",
		struct {
			Email  string
			Status string
		}{
			Email:  request.Email,
			Status: "ACTIVE",
		},
	)

	return result.SuccessWithValue(200, request)
}

// Send the user an OTP via email
// TODO - add a TTL to the verification code
func SendEmailChallenge(res result.ResultValue) *result.Result {
	request := res.(LoginChallengeRequest)

	response := container.SesClient().
		BuildEmailRequest(request.Email, request.Code).
		Send()

	if !response.IsSuccess {
		logger.Error(
			"Failed to send login challenge email",
			struct {
				Email string
				Error coreTypes.PasswordCaddyError
			}{
				Email: request.Email,
				Error: response.Error,
			},
		)

		return result.Failure(
			response.Error.StatusCode,
			response.Error.Message,
		)
	}

	logger.Info(
		"Successfully sent challenge email",
		struct {
			Email     string
			MessageId string
		}{
			Email:     request.Email,
			MessageId: response.Data.(string),
		},
	)

	return result.SuccessWithValue(202, request)
}

// Save the OTP in DynamoDB for verification use later
func AddOTPToDynamo(res result.ResultValue) *result.Result {
	request := res.(LoginChallengeRequest)

	dynamoRequest := dynamoclient.DyanamoUpdateRequest{
		Key: request.Email,
		Values: map[string]dynamoclient.DynamoUpdateItem{
			"VERIFICATION_CODE": {
				Action: types.AttributeActionPut,
				Value:  request.Code,
			},
		},
	}

	response := container.DynamoClient().
		Update(dynamoRequest)

	if !response.IsSuccess {
		logger.Error(
			"Failed to save OTP to DynamoDB",
			struct {
				Email string
				Error coreTypes.PasswordCaddyError
			}{
				Email: request.Email,
				Error: response.Error,
			},
		)

		return result.Failure(
			response.Error.StatusCode,
			response.Error.Message,
		)
	}

	logger.Info(
		"Saved OTP to DynamoDB",
		struct{ Email string }{
			Email: request.Email,
		},
	)

	return result.Success(202)
}

// Handle the login challenge request
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Init(event).
		Then(GetEmailStatus).
		Then(UpdateEmailStatusInDynamo).
		Then(SendEmailChallenge).
		Then(AddOTPToDynamo).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

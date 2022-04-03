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

// Send the user an OTP via email
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
		struct{ Email string }{
			Email: request.Email,
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
		Then(SendEmailChallenge).
		Then(AddOTPToDynamo).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

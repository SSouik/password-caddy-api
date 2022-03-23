package main

import (
	"errors"
	"password-caddy/api/src/core/container"
	"password-caddy/api/src/lib/dynamoclient"
	"password-caddy/api/src/lib/result"
	"password-caddy/api/src/lib/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LoginChallengeRequest struct {
	Email string
	Code  string
}

// Initialize the Login Challenge Request request
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
		return result.Failure(
			response.Error.StatusCode,
			errors.New(response.Error.Message),
		)
	}

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
		return result.Failure(
			response.Error.StatusCode,
			errors.New(response.Error.Message),
		)
	}

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

package main

import (
	"password-caddy/api/core/container"
	"password-caddy/api/core/types"
	"password-caddy/api/lib/dynamoclient"
	"password-caddy/api/lib/logger"
	"password-caddy/api/lib/result"
	"password-caddy/api/lib/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginVerificationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func Init(event events.APIGatewayProxyRequest) *result.Result {
	var request LoginVerificationRequest

	err := util.DeserializeJson(event.Body, &request)

	if err != nil {
		return result.Failure(500, err.Error())
	}

	request.Email = event.PathParameters["email"]

	return result.SuccessWithValue(200, request)
}

// Check to see if the request code matches the one stored
func VerifyCode(res result.ResultValue) *result.Result {
	request := res.(LoginVerificationRequest)

	dynamoRequest := dynamoclient.DynamoGetRequest{
		Key: request.Email,
	}

	response := container.DynamoClient().
		Get(dynamoRequest).
		AsUser()

	if !response.IsSuccess {
		logger.Error(
			"Failed to fetch user data",
			struct {
				Email string
				Error types.PasswordCaddyError
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

	user := response.Data.(types.PasswordCaddyUser)

	if user.VerificationCode.Value != request.Code {
		logger.Warn(
			"Requested verification code does not match one on record",
			struct{ Email string }{
				Email: user.UserId.Value,
			},
		)

		return result.Failure(401, "Unauthorized login attempt")
	}

	logger.Info(
		"Successfully verified login code",
		struct{ Email string }{
			Email: request.Email,
		},
	)

	return result.Success(201)
}

// Handle the login verification request
// TODO - Create JWT Auth token and send to in response
func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Init(event).
		Then(VerifyCode).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

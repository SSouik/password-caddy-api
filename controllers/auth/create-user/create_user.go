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

type CreateUserRequest struct {
	Email string `json:"userId"`
}

// Initialize the Create User Request
func Init(event events.APIGatewayProxyRequest) *result.Result {
	var request CreateUserRequest

	err := util.DeserializeJson(event.Body, &request)

	if err != nil {
		return result.Failure(500, err.Error())
	}

	return result.SuccessWithValue(200, request)
}

// Check if the email requested is already is use and active
func CheckIfUserAlreadyExists(res result.ResultValue) *result.Result {
	request := res.(CreateUserRequest)

	dynamoRequest := dynamoclient.DynamoGetRequest{
		Key: request.Email,
	}

	response := container.DynamoClient().
		Get(dynamoRequest).
		AsUser()

	if !response.IsSuccess {
		logger.Error(
			"Failed to check if email already exists and is active",
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

	// If the requested email is associated with an active user,
	// fail and do not create another record
	if user.Status.Value == "ACTIVE" {
		logger.Warn(
			"Failed to create user since email already exists and is active",
			struct{ Email string }{
				Email: user.UserId.Value,
			},
		)

		return result.Failure(409, "Email is already active")
	}

	logger.Info(
		"Requested email to create account is acceptable",
		struct{ Email string }{
			Email: request.Email,
		},
	)

	return result.SuccessWithValue(202, request)
}

// Create a new record in DynamoDB with the request email
func CreateUser(res result.ResultValue) *result.Result {
	request := res.(CreateUserRequest)

	dynamoRequest := dynamoclient.DynamoPutRequest{
		Key: request.Email,
		Values: map[string]string{
			"STATUS": "PENDING_REGISTRATION",
		},
	}

	response := container.DynamoClient().
		Put(dynamoRequest)

	if !response.IsSuccess {
		logger.Error(
			"Failed to create a user",
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

	logger.Info(
		"Successfully created a new user",
		struct {
			Email string
		}{
			Email: request.Email,
		},
	)

	return result.SuccessWithValue(201, request)
}

func SendVerificationEmail(res result.ResultValue) *result.Result {
	request := res.(CreateUserRequest)

	response := container.SesClient().
		SendVerificationEmail(request.Email)

	if !response.IsSuccess {
		logger.Error(
			"Failed to send verification email",
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

	logger.Info(
		"Verification email sent",
		struct {
			Email string
		}{
			Email: request.Email,
		},
	)

	return result.Success(201)
}

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Init(event).
		Then(CheckIfUserAlreadyExists).
		Then(CreateUser).
		Then(SendVerificationEmail).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

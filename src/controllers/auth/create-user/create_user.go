package main

import (
	"errors"

	"password-caddy/api/src/core/container"
	"password-caddy/api/src/core/types"
	"password-caddy/api/src/lib/dynamoclient"
	"password-caddy/api/src/lib/result"
	"password-caddy/api/src/lib/util"

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
		return result.Failure(500, err)
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
		return result.Failure(
			response.Error.StatusCode,
			response.Error,
		)
	}

	user := response.Data.(types.PasswordCaddyUser)

	// If the requested email is associated with an active user,
	// fail and do not create another record
	if user.Status.Value == "ACTIVE" {
		return result.Failure(409, errors.New("Email is already active"))
	}

	return result.SuccessWithValue(202, request)
}

// Create a new record in DynamoDB with the request email
// TODO - Send verification email via SES
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
		return result.Failure(
			response.Error.StatusCode,
			response.Error,
		)
	}

	return result.Success(201)
}

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Init(event).
		Then(CheckIfUserAlreadyExists).
		Then(CreateUser).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

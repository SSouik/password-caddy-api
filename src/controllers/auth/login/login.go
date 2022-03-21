package main

import (
	"errors"
	"password-caddy/api/src/core/container"
	"password-caddy/api/src/lib/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginResponse struct {
	Token      string `json:"token"`
	Expiration int64  `json:"exp"`
}

func GetToken() *result.Result {
	response := container.SesClient().
		BuildEmailRequest("samuel.souik@gmail.com").
		Send()

	if response.IsSuccess {
		return result.SuccessWithValue(200, LoginResponse{Token: response.MessageId})
	}

	return result.Failure(401, errors.New(response.ErrorMessage))
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return GetToken().ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

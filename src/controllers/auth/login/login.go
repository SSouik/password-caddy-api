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
		return result.SuccessWithValue(LoginResponse{Token: response.MessageId})
	}

	return result.Failure(401, errors.New(response.ErrorMessage))
	// var response LoginResponse
	// response.Token = config.Get("TEST_TOKEN", "default_token").ToString()
	// response.Expiration = config.Get("EXPIRATION", "1000").ToInt64()

	// return result.SuccessWithValue(response)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return result.Create().
		ThenApply(GetToken).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

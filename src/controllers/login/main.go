package main

import (
	"password-caddy/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func GetToken() *result.Result {
	var response LoginResponse
	response.Token = "some_token"

	return result.SuccessWithValue(response)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return result.Create().
		ThenApply(GetToken).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

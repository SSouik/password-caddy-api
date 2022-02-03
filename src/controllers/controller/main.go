package main

import (
	"password-caddy/api/src/lib/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func SayHello() *result.Result {
	var response Response
	response.Message = "Hello from Password Caddy Api!"

	return result.SuccessWithValue(response)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return result.Create().
		ThenApply(SayHello).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

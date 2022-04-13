package main

import (
	"password-caddy/api/lib/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type HealthCheckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func Init() *result.Result {
	request := HealthCheckResponse{
		Status:  200,
		Message: "Password Caddy is up and running",
		Version: "0.0.9",
	}

	return result.SuccessWithValue(200, request)
}

func Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Init().ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

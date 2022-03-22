package main

import (
	"errors"
	"password-caddy/api/src/core/container"
	"password-caddy/api/src/lib/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func SendEmailChallenge(email string) *result.Result {
	response := container.SesClient().
		BuildEmailRequest(email).
		Send()

	if !response.IsSuccess {
		return result.Failure(
			response.Error.StatusCode,
			errors.New(response.Error.Message),
		)
	}

	return result.Success(202)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return SendEmailChallenge(request.PathParameters["email"]).ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

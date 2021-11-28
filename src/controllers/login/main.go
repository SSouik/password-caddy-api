package main

import (
	"password-caddy/response"
	"password-caddy/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var responseBody LoginResponse
	responseBody.Token = "some_token"
	responseJson := util.SerializeJson(responseBody)

	return response.Create().
		WithStatus(200).
		WithBody(responseJson).
		ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

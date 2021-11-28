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

	resp := response.Create()

	responseBody.Token = "some_token"

	responseJson := util.SerializeJson(responseBody)

	response.WithStatus(200, &resp)
	response.WithBody(responseJson, &resp)

	return response.ToAPIGatewayResponse(resp), nil
}

func main() {
	lambda.Start(Handler)
}

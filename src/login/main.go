package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Token string `json:"token"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Response

	body.Token = "some_token"

	responseBody, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

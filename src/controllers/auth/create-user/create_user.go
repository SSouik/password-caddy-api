package main

import (
	"context"
	"errors"

	"password-caddy/api/src/core/container"
	"password-caddy/api/src/lib/dynamoclient"
	"password-caddy/api/src/lib/result"
	"password-caddy/api/src/lib/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

type Request struct {
	Email string `json:"userId"`
}

func GetRequest(event events.APIGatewayProxyRequest) (Request, error) {
	var request Request

	err := util.DeserializeJson(event.Body, &request)

	if err != nil {
		return request, err
	}

	return request, nil
}

func CreateUser(request Request) *result.Result {
	var dynamoRequest dynamoclient.DynamoPutRequest

	dynamoRequest.Key = request.Email
	dynamoRequest.Values = map[string]string{
		"STATUS": "PENDING_REGISTRATION",
	}

	response := container.DynamoClient().
		Put(dynamoRequest)

	if !response.IsSuccess {
		return result.Failure(
			response.Error.StatusCode,
			errors.New(response.Error.Message),
		)
	}

	return result.Success(201)
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request, err := GetRequest(event)

	if err != nil {
		return result.Failure(500, err).
			ToAPIGatewayResponse()
	}

	return CreateUser(request).ToAPIGatewayResponse()
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"context"
	"encoding/json"

	"password-caddy/api/src/lib/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Response struct {
	Message string `json:"message"`
}

type Request struct {
	UserId string `json:"userId"`
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Response

	var request Request

	jsonErr := util.DeserializeJson(event.Body, &request)

	if jsonErr != nil {
		body.Message = jsonErr.Error()
		res, _ := json.Marshal(body)

		return events.APIGatewayProxyResponse{
			Body:       string(res),
			StatusCode: 500,
		}, nil
	}

	cfg, cfgerr := config.LoadDefaultConfig(context.TODO())

	if cfgerr != nil {
		body.Message = cfgerr.Error()
		res, _ := json.Marshal(body)

		return events.APIGatewayProxyResponse{
			Body:       string(res),
			StatusCode: 500,
		}, nil
	}

	// This works and will create a user in the data base
	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("password-caddy-dev"),
		Item: map[string]types.AttributeValue{
			"USER_ID": &types.AttributeValueMemberS{
				Value: request.UserId,
			},
		},
	}

	response, err := client.PutItem(ctx, input)

	if err != nil {
		body.Message = err.Error()
		res, _ := json.Marshal(body)

		return events.APIGatewayProxyResponse{
			Body:       string(res),
			StatusCode: 500,
		}, nil
	}

	responseBody := util.SerializeJson(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

package dynamoclient

import (
	"context"
	"password-caddy/api/src/core/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClient struct {
	client    *dynamodb.Client
	tableName string
}

type DynamoResponse struct {
	IsSuccess    bool
	ErrorMessage string
}

/*
Create a new instance of the AWS DynamoDB Client
*/
func Create(awsConfig aws.Config) *DynamoClient {
	var dynamo DynamoClient

	dynamo.client = dynamodb.NewFromConfig(awsConfig)

	// Get the table name from the config
	dynamo.tableName = config.Get("DYNAMO_TABLE", "password-caddy-dev").
		ToString()

	return &dynamo
}

/*
Map a string -> string JSON object to a DynamoDB Item
Update if additional types are needed
*/
func ConvertToDynamoItem(obj map[string]string) map[string]types.AttributeValue {
	dynamoItem := make(map[string]types.AttributeValue)

	for key, value := range obj {
		dynamoItem[key] = &types.AttributeValueMemberS{
			Value: value,
		}
	}

	return dynamoItem
}

/*
Put an item in the DynamoDB table

@see - https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.PutItem
*/
func (dynamo *DynamoClient) Put(item map[string]string) *DynamoResponse {
	var putInput *dynamodb.PutItemInput

	putInput = &dynamodb.PutItemInput{
		TableName: aws.String(dynamo.tableName),
		Item:      ConvertToDynamoItem(item),
	}

	_, err := dynamo.client.PutItem(context.TODO(), putInput)

	if err != nil {
		return Failure(err.Error())
	}

	return Success()
}

/*
Create A successful Dynamo response
*/
func Success() *DynamoResponse {
	return &DynamoResponse{
		IsSuccess: true,
	}
}

/*
Create a failure Dynamo response
*/
func Failure(errorMessage string) *DynamoResponse {
	return &DynamoResponse{
		IsSuccess:    false,
		ErrorMessage: errorMessage,
	}
}

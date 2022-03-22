package dynamoclient

import (
	"context"
	"errors"

	apiError "password-caddy/api/src/core/passwordcaddyerror"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoClient struct {
	Client *dynamodb.Client
	Config DynamoConfig
}

type DynamoConfig struct {
	TableName string
}

type DynamoResponse struct {
	IsSuccess bool
	Error     apiError.PasswordCaddyError
}

/*
Create a new instance of the AWS DynamoDB Client
*/
func Create(awsConfig aws.Config) *DynamoClient {
	var dynamo DynamoClient
	dynamo.Client = dynamodb.NewFromConfig(awsConfig)

	return &dynamo
}

func (dynamo *DynamoClient) WithConfig(config DynamoConfig) *DynamoClient {
	dynamo.Config = config
	return dynamo
}

/*
Put an item in the DynamoDB table

@see - https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.PutItem
*/
func (dynamo *DynamoClient) Put(item map[string]string) *DynamoResponse {
	var putInput *dynamodb.PutItemInput

	putInput = &dynamodb.PutItemInput{
		TableName: aws.String(dynamo.Config.TableName),
		Item:      ConvertToDynamoItem(item),
	}

	_, err := dynamo.Client.PutItem(context.TODO(), putInput)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(apiError.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiError.PasswordCaddyError{
			StatusCode: 504,
			Message:    err.Error(),
		})
	}

	return Success()
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
func Failure(pcError apiError.PasswordCaddyError) *DynamoResponse {
	return &DynamoResponse{
		IsSuccess: false,
		Error:     pcError,
	}
}

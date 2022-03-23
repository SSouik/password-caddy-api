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

type DynamoPutRequest struct {
	Key    string
	Values map[string]string
}

type DyanamoUpdateRequest struct {
	Key    string
	Values map[string]DynamoUpdateItem
}

type DynamoUpdateItem struct {
	Action types.AttributeAction
	Value  string
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

// Put an item in the DynamoDB table. Creates a new item and overwrites if it exists
//
// @see - https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.PutItem
func (dynamo *DynamoClient) Put(request DynamoPutRequest) *DynamoResponse {
	var putInput *dynamodb.PutItemInput

	request.Values["USER_ID"] = request.Key

	putInput = &dynamodb.PutItemInput{
		TableName: aws.String(dynamo.Config.TableName),
		Item:      ConvertToDynamoPutItem(request.Values),
	}

	_, err := dynamo.Client.PutItem(context.TODO(), putInput)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(apiError.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiError.PasswordCaddyError{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	return Success()
}

// Update an Existing item in DynamoDB
//
// @see - https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.UpdateItem
func (dynamo *DynamoClient) Update(request DyanamoUpdateRequest) *DynamoResponse {
	var updateInput *dynamodb.UpdateItemInput

	updateInput = &dynamodb.UpdateItemInput{
		TableName: aws.String(dynamo.Config.TableName),
		Key: map[string]types.AttributeValue{
			"USER_ID": &types.AttributeValueMemberS{
				Value: request.Key,
			},
		},
		AttributeUpdates: ConvertToDynamoUpdateItem(request.Values),
	}

	_, err := dynamo.Client.UpdateItem(context.TODO(), updateInput)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(apiError.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiError.PasswordCaddyError{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	return Success()
}

/*
Map a string -> string JSON object to a DynamoDB PutItem
Update if additional types are needed
*/
func ConvertToDynamoPutItem(obj map[string]string) map[string]types.AttributeValue {
	dynamoItem := make(map[string]types.AttributeValue)

	for key, value := range obj {
		dynamoItem[key] = &types.AttributeValueMemberS{
			Value: value,
		}
	}

	return dynamoItem
}

/*
Map a string -> string JSON object to a DynamoDB UpdateItem
Update if additional types are needed
*/
func ConvertToDynamoUpdateItem(item map[string]DynamoUpdateItem) map[string]types.AttributeValueUpdate {
	dynamoItem := make(map[string]types.AttributeValueUpdate)

	for key, value := range item {
		dynamoItem[key] = types.AttributeValueUpdate{
			Action: value.Action,
			Value: &types.AttributeValueMemberS{
				Value: value.Value,
			},
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

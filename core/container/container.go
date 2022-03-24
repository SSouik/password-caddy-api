package container

import (
	"context"
	"fmt"
	appConfig "password-caddy/api/core/config"
	"password-caddy/api/lib/dynamoclient"
	"password-caddy/api/lib/sesclient"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func TestGetConfig() string {
	return appConfig.Get("TEST_TOKEN", "FOOBAR").ToString()
}

/*
Load the Default AWS config with AWS credentials
*/
func LoadAwsConfig() aws.Config {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())

	// If the config couldn't load, then panic and fail
	if err != nil {
		panic(fmt.Sprint("FATAL: Failed to load AWS config!"))
	}

	return cfg
}

func SesClient() *sesclient.SesClient {
	return sesclient.Create(LoadAwsConfig())
}

func DynamoClient() *dynamoclient.DynamoClient {
	var config dynamoclient.DynamoConfig

	config = dynamoclient.DynamoConfig{
		TableName: appConfig.Get("DYNAMO_TABLE", "password-caddy-dev").ToString(),
	}

	return dynamoclient.Create(LoadAwsConfig()).
		WithConfig(config)
}

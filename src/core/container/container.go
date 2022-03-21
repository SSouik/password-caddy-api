package container

import (
	"context"
	"fmt"
	"password-caddy/api/src/core/config"
	"password-caddy/api/src/lib/sesclient"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func TestGetConfig() string {
	return config.Get("TEST_TOKEN", "FOOBAR").ToString()
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

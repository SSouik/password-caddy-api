# This script will bootstrap a Hello world controller written in golang
# at the location you specify relative to the src directory
# Example: ./scripts/CreateController.ps1 foo
# The above example will create a directory named "foo" at src/foo
# and bootstrap a go module called "foo" that is set up to run as
# an AWS Lambda function

$ControllerName = $args[0]

if ($null -eq $ControllerName)
{
    $ControllerName = "controller"
}

$Dir = (Get-Location).Path
$Path = "$($Dir)\src\controllers\$($ControllerName)";

$ControllerTemplate = 
'package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body Response

	body.Message = "Hello world!"

	responseBody, _ := json.Marshal(body)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}'

Write-Host "Checking if directory src/controllers/$($ControllerName) exists..."

if (Test-Path -Path "$($Path)")
{
    Write-Host "Error: src/controllers/$($ControllerName) already exists"
    exit(1)
}

Write-Host "Info: Creating directory src/controllers/$($ControllerName) ..."

New-Item -Path "$($Path)" -ItemType Directory | Out-Null

Write-Host "Info: Creating file src/controllers/$($ControllerName)/main.go ..."

New-Item -Path "$($Path)\main.go" -ItemType File | Out-Null
Set-Content -Path "$($Path)\main.go" -Value $ControllerTemplate | Out-Null

Write-Host "Info: Initializing module..."
Set-Location -Path $Path

# Set up go module with defaults
go mod init "password-caddy/$($ControllerName)"
go mod tidy
go get github.com/aws/aws-lambda-go
go install

Set-Location -Path $Dir

Write-Host "----------"
Write-Host "Controller created. Make sure to add the controller to template.yml"
Write-Host ""

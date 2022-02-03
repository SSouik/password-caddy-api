# This script will bootstrap a Hello world controller written in golang
# at the location you specify relative to the src/controllers directory
# Example: ./scripts/CreateController.ps1 foo
# The above example will create a directory named "foo" at src/controllers/foo
# and a file main.go set up to be run as an AWS Lambda function

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
	"password-caddy/api/src/lib/result"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func SayHello() *result.Result {
	var response Response
	response.Message = "Hello from Password Caddy Api!"

	return result.SuccessWithValue(response)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return result.Create().
		ThenApply(SayHello).
		ToAPIGatewayResponse()
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

Write-Host "----------"
Write-Host "Controller created. Make sure to add the controller to template.yml"
Write-Host ""

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

	"password-caddy/result"

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

Write-Host "Info: Initializing module..."
Set-Location -Path $Path

# Set up go module with defaults
go mod init "password-caddy/$($ControllerName)"
go mod tidy
go get github.com/aws/aws-lambda-go
go install

Add-Content -Path "$($Path)\go.mod" -Value "require password-caddy/result v1.0.0" | Out-Null
Add-Content -Path "$($Path)\go.mod" -Value "" | Out-Null
Add-Content -Path "$($Path)\go.mod" -Value "replace password-caddy/result => ../../lib/result" | Out-Null

Set-Location -Path $Dir

Write-Host "----------"
Write-Host "Controller created. Make sure to add the controller to template.yml"
Write-Host ""

# Password Caddy API
Password Caddy API built using AWS SAM and GoLang. This API leverages AWS SAM to create a serverless architecture API utilizing, API Gateway, Lambda, and Dynamo DB.

<br/>

## Table of Contents
* [File Structure](#file-structure)
* [Requirements](#requirements)
* [Getting Started](#getting-started)
    * [Install Modules](#install-modules)
    * [Validate](#validate)
    * [Build](#build)
    * [Local Start](#local-start)
    * [Unit Tests](#unit-tests)

<br/>

## File Structure
```bash
.
├── Makefile                               <-- Make to automate build
├── README.md                              <-- This instructions file
├── scripts                                <-- Contains useful scripts for local development and CI
│   ├── CreateController.ps1               <-- Powershell script to bootstrap a controller under the src/controllers dir
│   ├── set_env.sh                         <-- Changes environment variables in template.yaml for local development
│   └── version_check                      <-- Checks if current version is greater than upstream version
├── src                                    <-- Source code for a lambda function
│   ├── controllers                        <-- Contains all the controllers (Lambda) code
│   │   └── auth                           <-- Contains all the Auth controllers (i.e login, registration)
│   │       ├── create-user                <-- The controller folder for the Lambda function to create a user
│   │       │    ├── create_user.go        <-- Create User Lambda function code
│   │       │    └── create_user_test.go   <-- Unit tests for create user Lambda function
│   │       └── login                      <-- The controller folder for the Lambda function to login
│   │           ├── login.go               <-- Login Lambda function code
│   │           └── login_test.go          <-- Unit tests for login Lambda function
│   ├── core                               <-- Contains core packages (Types, config, etc)
│   │   └── config                         <-- Config module (Contains the config core logic)
│   │       ├── config.go                  <-- Config module (Hanldes fetching environment variables and converting them to int, bool, etc)
│   │       └── config_test.go             <-- Unit tests for config
│   └── lib                                <-- Contains all the controllers (Lambda) code
│       ├── result                         <-- The result module (Contains the types and functions for a global result object)
│       │   ├── result.go                  <-- The result code
│       │   └── result_test.go             <-- Unit tests for result
│       └── util                           <-- The util module (Contains static utility methods)
│           ├── util.go                    <-- Collection of utility functions
│           └── util_test.go               <-- Unit tests for util
├── go.mod                                 <-- Root go module file
├── go.sum                                 <-- Root go module sum file
├── sam.version                            <-- SAM version file (Just shows the current version of SAM used)
├── template.local.yml                     <-- SAM template file for local development
├── template.yml                           <-- SAM template file for AWS deployments
└── version.json                           <-- Version of the API in format major.minor.patch (i.e 1.0.0)
```

<br/>

## Requirements

* [AWS CLI](https://aws.amazon.com/cli/)
* [Docker](https://www.docker.com/community-edition)
* [Golang](https://golang.org) version `1.17`
* [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) version `1.35`

<br/>

## Getting Started
This project has a `Makefile`. To install on Windows, use `chocolately`
```powershell
choco install make
```

### Install Modules
```bash
go mod vendor
```
> This will install go modules in `go.mod` into the `vendor` directory

<br/>

### Validate
Validate `template.yml`
```bash
make validate
```

Validate `template.local.yml`
```bash
make validate-local
```

<br/>

### Build
Build `template.yml`
```bash
make build
```

Build `template.local.yml`
```bash
make build-local
```

<br/>

### Local Start
Start the API locally. Requires Docker to be running.
```bash
make api
```

<br/>

### Unit Tests
Run unit tests
```bash
make test
```

Run unit tests with coverage report. Opens the report afterwards.
```bash
make coverage
```

<br/>

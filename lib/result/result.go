package result

import (
	"password-caddy/api/core/types"
	"password-caddy/api/lib/util"

	"github.com/aws/aws-lambda-go/events"
)

type ResultError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type ResultValue interface{}

type Resulter interface {
	GetValue() ResultValue
	Then(f func(res ResultValue) *Result) *Result
}

type Result struct {
	IsSuccess  bool
	StatusCode int
	Value      ResultValue
	Error      types.PasswordCaddyError
}

func (result *Result) GetValue() ResultValue {
	return result.Value
}

func (result *Result) Then(f func(res ResultValue) *Result) *Result {
	if !result.IsSuccess {
		return result
	}

	return f(result.GetValue())
}

// Create a new successful Result
func Success(statusCode int) *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: statusCode,
		Value:      nil,
		Error:      *new(types.PasswordCaddyError),
	}
}

// Create a new successful Result with a value
func SuccessWithValue(statusCode int, value ResultValue) *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: statusCode,
		Value:      value,
		Error:      *new(types.PasswordCaddyError),
	}
}

// Create a new failure Result with a status code and error
func Failure(statusCode int, message string) *Result {
	return &Result{
		IsSuccess:  false,
		StatusCode: statusCode,
		Value:      nil,
		Error: types.PasswordCaddyError{
			StatusCode: statusCode,
			Message:    message,
		},
	}
}

// Convert the Result to an API Gateway Proxy Reponse
func (result *Result) ToAPIGatewayResponse() (events.APIGatewayProxyResponse, error) {
	defaultHeaders := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: result.StatusCode,
		Headers:    defaultHeaders,
	}

	if result.IsSuccess && result.Value == nil {
		return response, nil
	} else if result.IsSuccess {
		body := util.SerializeJson(result.Value)
		response.Body = body
	} else {
		apiResponse := types.PasswordCaddyErrorResponse{
			Error: ResultError{
				StatusCode: result.Error.StatusCode,
				Message:    result.Error.Message,
			},
		}

		body := util.SerializeJson(apiResponse)

		response.Body = body
	}

	return response, nil
}

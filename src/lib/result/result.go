package result

import (
	"password-caddy/api/src/lib/util"

	"github.com/aws/aws-lambda-go/events"
)

type ResultError struct {
	Message string `json:"message"`
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
	Error      error
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
		Error:      nil,
	}
}

// Create a new successful Result with a value
func SuccessWithValue(statusCode int, value ResultValue) *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: statusCode,
		Value:      value,
		Error:      nil,
	}
}

// Create a new failure Result with a status code and error
func Failure(statusCode int, err error) *Result {
	return &Result{
		IsSuccess:  false,
		StatusCode: statusCode,
		Value:      nil,
		Error:      err,
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
	}

	if result.IsSuccess {
		body := util.SerializeJson(result.Value)
		response.Body = body
	} else {
		body := util.SerializeJson(ResultError{
			Message: result.Error.Error(),
		})

		response.Body = body
	}

	return response, nil
}

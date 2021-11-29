package result

import (
	"password-caddy/util"

	"github.com/aws/aws-lambda-go/events"
)

type ResultError struct {
	Message string `json:"message"`
}

type ResultValue interface{}

type Result struct {
	IsSuccess  bool
	StatusCode int
	Value      ResultValue
	Error      error
}

/*
  Create a new Result
*/
func Create() *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}
}

func CreateWithValue(value ResultValue) *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      value,
		Error:      nil,
	}
}

/********** SUCCESS FUNCS **********/

/*
  Create a new successful Result
*/
func Success() *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}
}

/*
  Create a new successful Result with a value
*/
func SuccessWithValue(value ResultValue) *Result {
	return &Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      value,
		Error:      nil,
	}
}

/********** END SUCCESS FUNCS **********/

/********** FAILURE FUNCS **********/

/*
  Create a new failure Result with a status code and error
*/
func Failure(statusCode int, err error) *Result {
	return &Result{
		IsSuccess:  false,
		StatusCode: statusCode,
		Value:      nil,
		Error:      err,
	}
}

/********** END FAILURE FUNCS **********/

/********** EXTENSION FUNCS **********/

/*
  Execute a function that returns a Result only when the
  Result is successful
*/
func (result *Result) ThenApply(fn func() *Result) *Result {
	if !result.IsSuccess {
		return Failure(result.StatusCode, result.Error)
	}

	return fn()
}

/*
  Apply a function to the Result's value
*/
func (result *Result) ThenApplyToValue(fn func(ResultValue) *Result) *Result {
	if !result.IsSuccess {
		return Failure(result.StatusCode, result.Error)
	}

	return fn(result.Value)
}

/*
  Apply a predicate function to the Result and if it returns true,
  then return a successful Result
*/
func (result *Result) SuccessWhen(fn func(*Result) bool) *Result {
	if fn(result) {
		return Success()
	}

	return Failure(result.StatusCode, result.Error)
}

/*
  Apply a predicate function to the Result's value and if it returns true,
  then return a successful Result
*/
func (result *Result) SuccessWhenValue(fn func(ResultValue) bool) *Result {
	if fn(result.Value) {
		return Success()
	}

	return Failure(result.StatusCode, result.Error)
}

/*
  Convert the Result to an API Gateway Proxy Reponse
*/
func (result *Result) ToAPIGatewayResponse() (events.APIGatewayProxyResponse, error) {
	defaultHeaders := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: result.StatusCode,
		Headers:    defaultHeaders,
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

/********** END EXTENSION FUNCS **********/

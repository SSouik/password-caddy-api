package result

import (
	"password-caddy/api/core/types"
	"password-caddy/api/lib/util"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func AddOne(res ResultValue) *Result {
	return SuccessWithValue(200, res.(int)+1)
}

func TestResultGetValue(t *testing.T) {
	result := SuccessWithValue(200, "abc")

	actual := result.GetValue().(string)
	expected := "abc"

	if actual != expected {
		t.Errorf("FAILED - TestResultGetValue | Actual: %s | Expected: %s", actual, expected)
	}
}

func TestResultThenMethodWithSuccessfulResult(t *testing.T) {
	result := SuccessWithValue(200, 1)

	actual := result.Then(AddOne).GetValue().(int)
	expected := 2

	if actual != expected {
		t.Errorf("FAILED - TestResultThenMethodWithSuccessfulResult | Actual: %d | Expected: %d", actual, expected)
	}
}

func TestResultThenMethodWithFailureResult(t *testing.T) {
	result := Failure(500, "Internal Server Error")

	actual := result.Then(AddOne).Error.Error()
	expected := "Internal Server Error"

	if actual != expected {
		t.Errorf("FAILED - TestResultThenMethodWithFailureResult | Actual: %s | Expected: %s", actual, expected)
	}
}

func TestSuccess(t *testing.T) {
	actual := Success(202)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 202,
		Value:      nil,
		Error:      *new(types.PasswordCaddyError),
	}

	if actual.IsSuccess != expected.IsSuccess {
		t.Errorf("FAILED - TestSuccess | Actual: %t | Expected: %t", actual.IsSuccess, expected.IsSuccess)
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestSuccess | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Error != expected.Error {
		t.Errorf("FAILED - TestSuccess | Actual: %+v | Expected: %+v", actual.Error, expected.Error)
	}
}

func TestSuccessWithValue(t *testing.T) {
	actual := SuccessWithValue(200, "Foo")
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      "Foo",
		Error:      *new(types.PasswordCaddyError),
	}

	if actual.IsSuccess != expected.IsSuccess {
		t.Errorf("FAILED - TestSuccessWithValue | Actual: %t | Expected: %t", actual.IsSuccess, expected.IsSuccess)
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestSuccessWithValue | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Error != expected.Error {
		t.Errorf("FAILED - TestSuccessWithValue | Actual: %+v | Expected: %+v", actual.Error, expected.Error)
	}
}

func TestFailure(t *testing.T) {
	actual := Failure(404, "Not Found")
	expected := Result{
		IsSuccess:  false,
		StatusCode: 404,
		Value:      nil,
		Error: types.PasswordCaddyError{
			StatusCode: 404,
			Message:    "Not Found",
		},
	}

	if actual.IsSuccess != expected.IsSuccess {
		t.Errorf("FAILED - TestFailure | Actual: %t | Expected: %t", actual.IsSuccess, expected.IsSuccess)
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestFailure | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Error != expected.Error {
		t.Errorf("FAILED - TestFailure | Actual: %+v | Expected: %+v", actual.Error, expected.Error)
	}
}

func TestToAPIGatewayResponseWithSuccessfulResult(t *testing.T) {
	res := SuccessWithValue(200, "Foo")
	actual, _ := res.ToAPIGatewayResponse()
	expected := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		Body: util.SerializeJson("Foo"),
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - StatusCode | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Headers["Accept"] != expected.Headers["Accept"] || actual.Headers["Content-Type"] != expected.Headers["Content-Type"] {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - Headers | Actual: %+v | Expected: %+v", actual.Headers, expected.Headers)
	}

	if actual.Body != expected.Body {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - Body | Actual: %s | Expected: %s", actual.Body, expected.Body)
	}
}

func TestToAPIGatewayResponseWithSuccessfulResultWithNoBody(t *testing.T) {
	res := Success(202)
	actual, _ := res.ToAPIGatewayResponse()
	expected := events.APIGatewayProxyResponse{
		StatusCode: 202,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - StatusCode | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Headers["Accept"] != expected.Headers["Accept"] || actual.Headers["Content-Type"] != expected.Headers["Content-Type"] {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - Headers | Actual: %+v | Expected: %+v", actual.Headers, expected.Headers)
	}

	if actual.Body != expected.Body {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithSuccessfulResult - Body | Actual: %s | Expected: %s", actual.Body, expected.Body)
	}
}

func TestToAPIGatewayResponseWithFailureResult(t *testing.T) {
	res := Failure(500, "Internal Error")
	actual, _ := res.ToAPIGatewayResponse()
	expected := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		Body: "{\"error\":{\"statusCode\":500,\"message\":\"Internal Error\"}}",
	}

	if actual.StatusCode != expected.StatusCode {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithFailureResult - StatusCode | Actual: %d | Expected: %d", actual.StatusCode, expected.StatusCode)
	}

	if actual.Headers["Accept"] != expected.Headers["Accept"] || actual.Headers["Content-Type"] != expected.Headers["Content-Type"] {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithFailureResult - Headers | Actual: %+v | Expected: %+v", actual.Headers, expected.Headers)
	}

	if actual.Body != expected.Body {
		t.Errorf("FAILED - TestToAPIGatewayResponseWithFailureResult - Body | Actual: %s | Expected: %s", actual.Body, expected.Body)
	}
}

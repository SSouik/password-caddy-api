package result

import (
	"errors"
	"password-caddy/api/src/lib/util"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestSuccess(t *testing.T) {
	actual := Success(202)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 202,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccess | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWithValue(t *testing.T) {
	actual := SuccessWithValue(200, "Foo")
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      "Foo",
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccessWithValue | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestFailure(t *testing.T) {
	err := errors.New("Not Found")

	actual := Failure(404, err)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 404,
		Value:      nil,
		Error:      err,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestFailure | Actual: %+v | Expected: %+v", *actual, expected)
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
	res := Failure(500, errors.New("Internal Error"))
	actual, _ := res.ToAPIGatewayResponse()
	expected := events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers: map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		Body: "{\"message\":\"Internal Error\"}",
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

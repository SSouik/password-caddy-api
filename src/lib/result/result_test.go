package result

import (
	"errors"
	"password-caddy/api/src/lib/util"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var INTERNAL_ERROR = errors.New("Internal Error")

func Add(x, y int) int {
	return x + y
}

func Concat(a, b string) string {
	return a + b
}

func SayHello() *Result {
	return SuccessWithValue("Hello")
}

func InternalError() *Result {
	return Failure(500, INTERNAL_ERROR)
}

func SayMessage(message ResultValue) *Result {
	return SuccessWithValue(message)
}

func InternalError2(message ResultValue) *Result {
	return Failure(500, INTERNAL_ERROR)
}

func PredicateTrue(result *Result) bool {
	return true
}

func PredicateFalse(result *Result) bool {
	return false
}

func PredicateTrue2(message ResultValue) bool {
	return true
}

func PredicateFalse2(message ResultValue) bool {
	return false
}

func TestCreateResult(t *testing.T) {
	actual := Create()
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestCreateResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestCreateResultWithValue(t *testing.T) {
	actual := CreateWithValue("Foo")
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      "Foo",
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestCreateResultWithValue | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccess(t *testing.T) {
	actual := Success()
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccess | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWithValue(t *testing.T) {
	actual := SuccessWithValue("Foo")
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

func TestThenApplyWithSuccessfulResult(t *testing.T) {
	res := Success()
	actual := res.ThenApply(SayHello)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      "Hello",
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyWithSuccessfulResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestThenApplyWhenFuncReturnsFailureResult(t *testing.T) {
	res := Failure(500, INTERNAL_ERROR)
	actual := res.ThenApply(SayHello)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 500,
		Value:      nil,
		Error:      INTERNAL_ERROR,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyWhenFuncReturnsFailureResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestThenApplyWithFailureResult(t *testing.T) {
	res := Success()
	actual := res.ThenApply(InternalError)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 500,
		Value:      nil,
		Error:      INTERNAL_ERROR,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyWithFailureResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestThenApplyToValueWithSuccessfulResult(t *testing.T) {
	res := SuccessWithValue("Foo")
	actual := res.ThenApplyToValue(SayMessage)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      "Foo",
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyToValueWithSuccessfulResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestThenApplyToValueWithFailureResult(t *testing.T) {
	res := Failure(500, INTERNAL_ERROR)
	actual := res.ThenApplyToValue(SayMessage)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 500,
		Value:      nil,
		Error:      INTERNAL_ERROR,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyToValueWithFailureResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestThenApplyToValueWhenFuncReturnsFailureResult(t *testing.T) {
	res := SuccessWithValue("Foo")
	actual := res.ThenApplyToValue(InternalError2)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 500,
		Value:      nil,
		Error:      INTERNAL_ERROR,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestThenApplyToValueWhenFuncReturnsFailureResult | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWhenPredicateReturnsTrue(t *testing.T) {
	res := Success()
	actual := res.SuccessWhen(PredicateTrue)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccessWhenPredicateReturnsTrue | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWhenPredicateReturnsFalse(t *testing.T) {
	res := Success()
	actual := res.SuccessWhen(PredicateFalse)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccessWhenPredicateReturnsFalse | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWhenValuePredicateReturnsTrue(t *testing.T) {
	res := SuccessWithValue("message")
	actual := res.SuccessWhenValue(PredicateTrue2)
	expected := Result{
		IsSuccess:  true,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccessWhenValuePredicateReturnsTrue | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestSuccessWhenValuePredicateReturnsFalse(t *testing.T) {
	res := SuccessWithValue("Foo")
	actual := res.SuccessWhenValue(PredicateFalse2)
	expected := Result{
		IsSuccess:  false,
		StatusCode: 200,
		Value:      nil,
		Error:      nil,
	}

	if *actual != expected {
		t.Errorf("FAILED - TestSuccessWhenValuePredicateReturnsFalse | Actual: %+v | Expected: %+v", *actual, expected)
	}
}

func TestToAPIGatewayResponseWithSuccessfulResult(t *testing.T) {
	res := SuccessWithValue("Foo")
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

func TestToAPIGatewayResponseWithFailureResult(t *testing.T) {
	res := Failure(500, INTERNAL_ERROR)
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

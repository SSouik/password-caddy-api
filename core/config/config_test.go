package config

import (
	"os"
	"testing"
)

func TestParseIntWithValidInt(t *testing.T) {
	actual := ParseInt("123")
	expected := 123

	if actual != int64(expected) {
		t.Errorf("FAILED - TestParseIntWithValidInt | Actual: %d | Expected: %d", actual, expected)
	}
}

func TestParseIntWithInvalidInt(t *testing.T) {
	actual := ParseInt("abc123")
	expected := 0

	if actual != int64(expected) {
		t.Errorf("FAILED - TestParseIntWithInvalidInt | Actual: %d | Expected: %d", actual, expected)
	}
}

func TestParseBoolWithTrue(t *testing.T) {
	trueVals := [6]string{"1", "t", "T", "TRUE", "true", "True"}

	for _, val := range trueVals {
		actual := ParseBool(val)
		if actual != true {
			t.Logf("FAILED - TestParseBoolWithTrue | Actual: %v | Expected: %v", actual, true)
			t.Fail()
		}
	}
}

func TestParseBoolWithFalse(t *testing.T) {
	falseVals := [6]string{"0", "f", "F", "FALSE", "false", "False"}

	for _, val := range falseVals {
		actual := ParseBool(val)
		if actual != false {
			t.Logf("FAILED - TestParseBoolWithTrue | Actual: %v | Expected: %v", actual, true)
			t.Fail()
		}
	}
}

func TestParseBoolWithInvalidString(t *testing.T) {
	actual := ParseBool("abc")
	expected := false

	if actual != expected {
		t.Errorf("FAILED - TestParseBoolWithInvalidString | Actual: %v | Expected: %v", actual, expected)
	}
}

func TestToInt64WithValidInt(t *testing.T) {
	val := ConfigValue{Value: "123"}
	actual := val.ToInt64()
	expected := 123

	if actual != int64(expected) {
		t.Errorf("FAILED - TestToInt64WithValidInt | Actual: %d | Expected: %d", actual, expected)
	}
}

func TestToInt64WithInvalidInt(t *testing.T) {
	val := ConfigValue{Value: "foo"}
	actual := val.ToInt64()
	expected := 0

	if actual != int64(expected) {
		t.Errorf("FAILED - TestToInt64WithInvalidInt | Actual: %d | Expected: %d", actual, expected)
	}
}

func TestToString(t *testing.T) {
	val := ConfigValue{Value: "123"}
	actual := val.ToString()
	expected := "123"

	if actual != expected {
		t.Errorf("FAILED - TestToString | Actual: %s | Expected: %s", actual, expected)
	}
}

func TestGetWithEnvVarThatExists(t *testing.T) {
	os.Setenv("FOO", "bar")

	actual := Get("FOO", "default").ToString()
	expected := "bar"

	if actual != expected {
		t.Errorf("FAILED - TestGetWithEnvVarThatExists | Actual: %s | Expected: %s", actual, expected)
	}
}

func TestGetWithEnvVarThatDoesNotExist(t *testing.T) {
	actual := Get("BAR", "default").ToString()
	expected := "default"

	if actual != expected {
		t.Errorf("FAILED - TestGetWithEnvVarThatDoesNotExist | Actual: %s | Expected: %s", actual, expected)
	}
}

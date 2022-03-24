package util

import "testing"

type SampleObj struct {
	Foo string `json:"foo"`
	Bar int64  `json:"bar"`
}

type SampleObj2 struct {
	Foo string
	Bar int64
}

type SampleObj3 struct {
	foo string
	bar string
}

/***** DeserializeJson *****/

func TestDeserializeJsonWhenObjectHasJsonFields(t *testing.T) {
	var actual SampleObj
	json := "{\"foo\": \"message\", \"bar\": 123}"

	err := DeserializeJson(json, &actual)

	if err != nil {
		t.Errorf("DeserializeJson returned error: %s", err.Error())
	}

	expected := SampleObj{
		Foo: "message",
		Bar: 123,
	}

	if actual != expected {
		t.Errorf("Actual: %+v | Expected: %+v", actual, expected)
	}
}

func TestDeserializeJsonWhenObjectDoesNotHaveJsonFields(t *testing.T) {
	var actual SampleObj2
	json := "{\"foo\": \"message\", \"bar\": 123}"

	err := DeserializeJson(json, &actual)

	if err != nil {
		t.Errorf("DeserializeJson returned error: %s", err.Error())
	}

	expected := SampleObj2{
		Foo: "message",
		Bar: 123,
	}

	if actual != expected {
		t.Errorf("Actual: %+v | Expected: %+v", actual, expected)
	}
}

func TestDeserializeJsonWithInvalidJson(t *testing.T) {
	var actual SampleObj
	json := "{\"foo\": \"message\", bar: 123}"

	err := DeserializeJson(json, &actual)

	if err == nil {
		t.Errorf("DeserializeJson did not return error")
	}
}

/***** End DeserializeJson *****/

/***** SerializeJson *****/

func TestSerializeJson(t *testing.T) {
	obj := SampleObj{
		Foo: "message",
		Bar: 123,
	}

	expected := "{\"foo\":\"message\",\"bar\":123}"

	actual := SerializeJson(obj)

	if actual == "" {
		t.Errorf("SerializeJson failed to serialize")
	}

	if actual != expected {
		t.Errorf("Actual: %s | Expected: %s", actual, expected)
	}
}

func TestSerializeJsonWIthInvalidObject(t *testing.T) {
	obj := map[string]interface{}{
		"foo": make(chan int),
	}

	actual := SerializeJson(obj)

	if actual != "" {
		t.Errorf("SerializeJson failed to serialize")
	}
}

func TestGenerateOTPWithZeroLength(t *testing.T) {
	otp, _ := GenerateOTP(6)

	if len(otp) != 6 {
		t.Errorf("FAILED | GenerateOTP | Expected 6 | Actual %d", len(otp))
	}
}

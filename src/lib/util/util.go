package util

import (
	"crypto/rand"
	"encoding/json"
)

func DeserializeJson(data string, toObj interface{}) error {
	return json.Unmarshal([]byte(data), &toObj)
}

func SerializeJson(obj interface{}) string {
	json, err := json.Marshal(obj)

	if err != nil {
		return ""
	}

	return string(json)
}

func GenerateOTP(length int) (string, error) {
	const otpChars = "1234567890"

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)

	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)

	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

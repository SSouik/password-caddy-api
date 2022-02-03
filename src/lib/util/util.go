package util

import (
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

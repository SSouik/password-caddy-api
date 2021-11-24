package util

import (
	"encoding/json"
)

func ToJsonString(obj interface{}) string {
	json, err := json.Marshal(obj)

	if err != nil {
		return ""
	}

	return string(json)
}

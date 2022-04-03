package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"password-caddy/api/lib/util"

	"github.com/google/uuid"
)

const (
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
	LOG_LEVEL_DEBUG = "debug"
)

type LoggerInterface interface {
	Info(message string, details interface{})
	Warn(message string, details interface{})
	Error(message string, details interface{})
	Debug(message string, details interface{})
}

type LogMessage struct {
	ID      string      `json:"ID"`
	Level   string      `json:"Level"`
	Message string      `json:"Message"`
	Details interface{} `json:"Details"`
}

func PrettyPrintJson(str string) string {
	var prettyJSON bytes.Buffer

	err := json.Indent(&prettyJSON, []byte(str), "", "    ")

	if err != nil {
		log.Printf("Failed to log message: %s", err.Error())
	}

	return prettyJSON.String()
}

func Log(message LogMessage) {
	logMessage := PrettyPrintJson(util.SerializeJson(message))
	log.Print(logMessage)
}

func Info(message string, details interface{}) {
	logMessage := LogMessage{
		ID:      uuid.New().String(),
		Level:   LOG_LEVEL_INFO,
		Message: message,
		Details: details,
	}

	Log(logMessage)
}

func Warn(message string, details interface{}) {
	logMessage := LogMessage{
		ID:      uuid.New().String(),
		Level:   LOG_LEVEL_WARN,
		Message: message,
		Details: details,
	}

	Log(logMessage)
}

func Error(message string, details interface{}) {
	logMessage := LogMessage{
		ID:      uuid.New().String(),
		Level:   LOG_LEVEL_ERROR,
		Message: message,
		Details: details,
	}

	Log(logMessage)
}

func Debug(message string, details interface{}) {
	logMessage := LogMessage{
		ID:      uuid.New().String(),
		Level:   LOG_LEVEL_DEBUG,
		Message: message,
		Details: details,
	}

	Log(logMessage)
}

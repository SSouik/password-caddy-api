package passwordcaddyerror

import (
	"github.com/aws/smithy-go"
)

type PasswordCaddyError struct {
	StatusCode int
	Message    string
}

var (
	AWS_ERRORS_TO_STATUS_CODES = map[string]int{
		"ValidationException": 400,
	}
)

func AWSErrorToPasswordCaddyError(awsErr smithy.APIError) PasswordCaddyError {
	statusCode, exists := AWS_ERRORS_TO_STATUS_CODES[awsErr.ErrorCode()]

	if !exists {
		statusCode = 500
	}
	return PasswordCaddyError{
		StatusCode: statusCode,
		Message:    awsErr.ErrorMessage(),
	}
}

package types

/***** Custom Error *****/

type PasswordCaddyError struct {
	StatusCode int
	Message    string
}

func (pcError PasswordCaddyError) Error() string {
	return pcError.Message
}

/***** Generic Wrapper Values *****/

type StringValue struct {
	Value string `json:"Value"`
}

type NumberValue struct {
	Value int `json:"Value"`
}

/***** API Types *****/

type PasswordCaddyUser struct {
	UserId           StringValue `json:"USER_ID"`
	Status           StringValue `json:"STATUS"`
	VerificationCode StringValue `json:"VERIFICATION_CODE"`
}

type PasswordCaddyErrorResponse struct {
	Error interface{} `json:"error"`
}

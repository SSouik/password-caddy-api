package types

type StringValue struct {
	Value string `json:"Value"`
}

type PasswordCaddyUser struct {
	UserId           StringValue `json:"USER_ID"`
	Status           StringValue `json:"STATUS"`
	VerificationCode StringValue `json:"VERIFICATION_CODE"`
}

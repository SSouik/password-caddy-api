package okta

type OktaUser struct {
	ID      string      `json:"id"`
	Status  string      `json:"status"`
	Created string      `json:"created"`
	Profile OktaProfile `json:"profile"`
}

type OktaProfile struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	MobilePhone     string `json:"mobilePhone"`
	PasswordCaddyId string `json:"passwordCaddyId"`
	Login           string `json:"login"`
	Email           string `json:"email"`
}

type OktaClientConfig struct {
	BaseUrl string
	ApiKey  string
}

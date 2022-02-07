package oktaclient

import (
	"password-caddy/api/src/core/okta"

	"github.com/go-resty/resty/v2"
)

var DEFAULT_HEADERS = map[string]string{
	"Content-Type": "application/json",
	"Accept":       "application/json",
}

type OktaClient struct {
	Client *resty.Client
}

type OktaRequest struct {
	Request resty.Request
}

/********** SETUP FUNCTIONS **********/

func Create(oktaConfig okta.OktaClientConfig) *OktaClient {
	return &OktaClient{
		Client: resty.New().
			SetBaseURL(oktaConfig.BaseUrl).
			SetAuthScheme("SSWS").
			SetAuthToken(oktaConfig.ApiKey).
			SetRetryCount(3),
	}
}

func (client *OktaClient) Request() *OktaRequest {
	var request resty.Request
	request = *client.Client.R()
	request.URL = client.Client.BaseURL

	return &OktaRequest{
		Request: *request.
			SetAuthScheme(client.Client.AuthScheme).
			SetAuthToken(client.Client.Token).
			SetHeaders(DEFAULT_HEADERS),
	}
}

func (client *OktaClient) RequestWithHeaders(headers map[string]string) *OktaRequest {
	var request resty.Request
	request = *client.Client.R()
	request.URL = client.Client.BaseURL

	return &OktaRequest{
		Request: *request.
			SetAuthScheme(client.Client.AuthScheme).
			SetAuthToken(client.Client.Token).
			SetHeaders(DEFAULT_HEADERS).
			SetHeaders(headers),
	}
}

/********** END SETUP FUNCTIONS **********/

/********** HTTP METHODS FUNCTIONS **********/

func (request *OktaRequest) Get(endpoint string) *OktaRequest {
	request.Request.Method = resty.MethodGet
	request.Request.URL = request.Request.URL + endpoint
	return request
}

func (request *OktaRequest) Post(endpoint string, body interface{}) *OktaRequest {
	request.Request.Method = resty.MethodPost
	request.Request.URL = request.Request.URL + endpoint
	request.Request.Body = body
	return request
}

func (request *OktaRequest) Put(endpoint string, body interface{}) *OktaRequest {
	request.Request.Method = resty.MethodPut
	request.Request.URL = request.Request.URL + endpoint
	request.Request.Body = body
	return request
}

/********** END HTTP METHODS FUNCTIONS **********/

/********** ROUTE SETTER FUNCTIONS **********/

func (request *OktaRequest) SetRouteParam(param, value string) *OktaRequest {
	var oktaRequest OktaRequest
	oktaRequest.Request = *request.Request.SetPathParam(param, value)
	return &oktaRequest
}

func (request *OktaRequest) SetQueryParam(param, value string) *OktaRequest {
	var oktaRequest OktaRequest
	oktaRequest.Request = *request.Request.SetQueryParam(param, value)
	return &oktaRequest
}

/********** END ROUTE SETTER FUNCTIONS **********/

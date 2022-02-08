package oktaclient

import (
	"password-caddy/api/src/core/okta"
	"testing"

	"github.com/go-resty/resty/v2"
)

func BootstrapRequest(url string) *OktaRequest {
	var request resty.Request
	request = *resty.New().R()
	request.URL = url

	return &OktaRequest{
		Request: *request.
			SetAuthScheme("SSWS").
			SetAuthToken("abc132").
			SetHeaders(map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			}),
	}
}

/********** END SETUP TESTS **********/

func TestOktaClientCreate(t *testing.T) {
	actual := Create(okta.OktaClientConfig{
		BaseUrl: "http://example.com",
		ApiKey:  "abc123",
	})

	expected := OktaClient{
		Client: resty.New().
			SetBaseURL("http://example.com").
			SetAuthScheme("SSWS").
			SetAuthToken("abc123").
			SetRetryCount(3),
	}

	if actual.Client.BaseURL != expected.Client.BaseURL {
		t.Errorf("FAILED - TestOktaClientCreate | Actual: %s | Expected: %s", actual.Client.BaseURL, expected.Client.BaseURL)
	}

	if actual.Client.AuthScheme != expected.Client.AuthScheme {
		t.Errorf("FAILED - TestOktaClientCreate | Actual: %s | Expected: %s", actual.Client.AuthScheme, expected.Client.AuthScheme)
	}

	if actual.Client.Token != expected.Client.Token {
		t.Errorf("FAILED - TestOktaClientCreate | Actual: %s | Expected: %s", actual.Client.Token, expected.Client.Token)
	}

	if actual.Client.RetryCount != expected.Client.RetryCount {
		t.Errorf("FAILED - TestOktaClientCreate | Actual: %d | Expected: %d", actual.Client.RetryCount, expected.Client.RetryCount)
	}
}

func TestOktaClientRequest(t *testing.T) {
	actual := Create(okta.OktaClientConfig{
		BaseUrl: "http://example.com",
		ApiKey:  "abc123",
	}).Request()

	oktaClient := OktaClient{
		Client: resty.New().
			SetBaseURL("http://example.com").
			SetAuthScheme("SSWS").
			SetAuthToken("abc123").
			SetRetryCount(3),
	}

	oktaRequest := oktaClient.Client.R()
	oktaRequest.URL = "http://example.com"

	expected := OktaRequest{
		Request: *oktaRequest.
			SetAuthScheme("SSWS").
			SetAuthToken("abc123").
			SetHeaders(map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			}),
	}

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientRequest | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.AuthScheme != expected.Request.AuthScheme {
		t.Errorf("FAILED - TestOktaClientRequest | Actual: %s | Expected: %s", actual.Request.AuthScheme, expected.Request.AuthScheme)
	}

	if actual.Request.Token != expected.Request.Token {
		t.Errorf("FAILED - TestOktaClientRequest | Actual: %s | Expected: %s", actual.Request.Token, expected.Request.Token)
	}

	if actual.Request.Header.Get("Content-Type") != expected.Request.Header.Get("Content-Type") {
		t.Errorf("FAILED - TestOktaClientRequest | Actual: %s | Expected: %s", actual.Request.Header.Get("Content-Type"), expected.Request.Header.Get("Content-Type"))
	}

	if actual.Request.Header.Get("Accept") != expected.Request.Header.Get("Accept") {
		t.Errorf("FAILED - TestOktaClientRequest | Actual: %s | Expected: %s", actual.Request.Header.Get("Accept"), expected.Request.Header.Get("Accept"))
	}
}

func TestOktaClientRequestWithHeaders(t *testing.T) {
	actual := Create(okta.OktaClientConfig{
		BaseUrl: "http://example.com",
		ApiKey:  "abc123",
	}).RequestWithHeaders(map[string]string{
		"Foo": "bar",
		"bar": "foo",
	})

	oktaClient := OktaClient{
		Client: resty.New().
			SetBaseURL("http://example.com").
			SetAuthScheme("SSWS").
			SetAuthToken("abc123").
			SetRetryCount(3),
	}

	oktaRequest := oktaClient.Client.R()
	oktaRequest.URL = "http://example.com"

	expected := OktaRequest{
		Request: *oktaRequest.
			SetAuthScheme("SSWS").
			SetAuthToken("abc123").
			SetHeaders(map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"Foo":          "bar",
				"bar":          "foo",
			}),
	}

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.AuthScheme != expected.Request.AuthScheme {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.AuthScheme, expected.Request.AuthScheme)
	}

	if actual.Request.Token != expected.Request.Token {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.Token, expected.Request.Token)
	}

	if actual.Request.Header.Get("Content-Type") != expected.Request.Header.Get("Content-Type") {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.Header.Get("Content-Type"), expected.Request.Header.Get("Content-Type"))
	}

	if actual.Request.Header.Get("Accept") != expected.Request.Header.Get("Accept") {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.Header.Get("Accept"), expected.Request.Header.Get("Accept"))
	}

	if actual.Request.Header.Get("Foo") != expected.Request.Header.Get("Foo") {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.Header.Get("Foo"), expected.Request.Header.Get("Foo"))
	}

	if actual.Request.Header.Get("bar") != expected.Request.Header.Get("bar") {
		t.Errorf("FAILED - TestOktaClientRequestWithHeaders | Actual: %s | Expected: %s", actual.Request.Header.Get("bar"), expected.Request.Header.Get("bar"))
	}
}

/********** END SETUP TESTS **********/

/********** HTTP METHODS FUNCTION TESTS **********/

func TestOktaClientGet(t *testing.T) {
	url := "http://foo.com"
	actual := BootstrapRequest(url).Get("/abc/123")

	expected := BootstrapRequest(url)
	expected.Request.Method = resty.MethodGet
	expected.Request.URL = url + "/abc/123"

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientGet | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.Method != expected.Request.Method {
		t.Errorf("FAILED - TestOktaClientGet | Actual: %s | Expected: %s", actual.Request.Method, expected.Request.Method)
	}
}

func TestOktaClientPostWithJsonEncodedBody(t *testing.T) {
	url := "http://foo.com"
	body := "{\"foo\": \"bar\"}"

	actual := BootstrapRequest(url).Post("/abc/123", body)

	expected := BootstrapRequest(url)
	expected.Request.Method = resty.MethodPost
	expected.Request.URL = url + "/abc/123"
	expected.Request.Body = body

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientPostWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.Method != expected.Request.Method {
		t.Errorf("FAILED - TestOktaClientPostWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.Method, expected.Request.Method)
	}

	if actual.Request.Body != expected.Request.Body {
		t.Errorf("FAILED - TestOktaClientPostWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.Body, expected.Request.Body)
	}
}

func TestOktaClientPostWithStructBody(t *testing.T) {
	type FooObj struct {
		Foo string `json:"foo"`
	}

	url := "http://foo.com"

	body := FooObj{
		Foo: "bar",
	}

	actual := BootstrapRequest(url).Post("/abc/123", body)

	expected := BootstrapRequest(url)
	expected.Request.Method = resty.MethodPost
	expected.Request.URL = url + "/abc/123"
	expected.Request.Body = body

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientPostWithStructBody | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.Method != expected.Request.Method {
		t.Errorf("FAILED - TestOktaClientPostWithStructBody | Actual: %s | Expected: %s", actual.Request.Method, expected.Request.Method)
	}

	if actual.Request.Body != expected.Request.Body {
		t.Errorf("FAILED - TestOktaClientPostWithStructBody | Actual: %+v | Expected: %+v", actual.Request.Body, expected.Request.Body)
	}
}

func TestOktaClientPutWithJsonEncodedBody(t *testing.T) {
	url := "http://foo.com"
	body := "{\"foo\": \"bar\"}"

	actual := BootstrapRequest(url).Put("/abc/123", body)

	expected := BootstrapRequest(url)
	expected.Request.Method = resty.MethodPut
	expected.Request.URL = url + "/abc/123"
	expected.Request.Body = body

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientPutWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.Method != expected.Request.Method {
		t.Errorf("FAILED - TestOktaClientPutWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.Method, expected.Request.Method)
	}

	if actual.Request.Body != expected.Request.Body {
		t.Errorf("FAILED - TestOktaClientPutWithJsonEncodedBody | Actual: %s | Expected: %s", actual.Request.Body, expected.Request.Body)
	}
}

func TestOktaClientPutWithStructBody(t *testing.T) {
	type FooObj struct {
		Foo string `json:"foo"`
	}

	url := "http://foo.com"

	body := FooObj{
		Foo: "bar",
	}

	actual := BootstrapRequest(url).Put("/abc/123", body)

	expected := BootstrapRequest(url)
	expected.Request.Method = resty.MethodPut
	expected.Request.URL = url + "/abc/123"
	expected.Request.Body = body

	if actual.Request.URL != expected.Request.URL {
		t.Errorf("FAILED - TestOktaClientPutWithStructBody | Actual: %s | Expected: %s", actual.Request.URL, expected.Request.URL)
	}

	if actual.Request.Method != expected.Request.Method {
		t.Errorf("FAILED - TestOktaClientPutWithStructBody | Actual: %s | Expected: %s", actual.Request.Method, expected.Request.Method)
	}

	if actual.Request.Body != expected.Request.Body {
		t.Errorf("FAILED - TestOktaClientPutWithStructBody | Actual: %+v | Expected: %+v", actual.Request.Body, expected.Request.Body)
	}
}

/********** END HTTP METHODS FUNCTION TESTS **********/

/********** ROUTE SETTER FUNCTION TESTS **********/

func TestOktaClientRouteParam(t *testing.T) {
	url := "https://foo.com"

	actual := BootstrapRequest(url).
		Get("/abc/{foo}").
		SetRouteParam("foo", "123")

	expected := BootstrapRequest(url)
	expected.Request.URL = url + "/abc/{foo}"
	expected.Request.PathParams = map[string]string{
		"foo": "123",
	}

	if actual.Request.PathParams["foo"] != expected.Request.PathParams["foo"] {
		t.Errorf("FAILED - TestOktaClientRouteParam | Actual: %s | Expected: %s", actual.Request.PathParams["foo"], expected.Request.PathParams["foo"])
	}
}

func TestOktaClientRouteParamWithMultipleParams(t *testing.T) {
	url := "https://foo.com"

	actual := BootstrapRequest(url).
		Get("/abc/{foo}/{bar}").
		SetRouteParam("foo", "123").
		SetRouteParam("bar", "456")

	expected := BootstrapRequest(url)
	expected.Request.URL = url + "/abc/{foo}/{bar}"
	expected.Request.PathParams = map[string]string{
		"foo": "123",
		"bar": "456",
	}

	if actual.Request.PathParams["foo"] != expected.Request.PathParams["foo"] {
		t.Errorf("FAILED - TestOktaClientRouteParamWithMultipleParams | Actual: %s | Expected: %s", actual.Request.PathParams["foo"], expected.Request.PathParams["foo"])
	}

	if actual.Request.PathParams["bar"] != expected.Request.PathParams["bar"] {
		t.Errorf("FAILED - TestOktaClientRouteParamWithMultipleParams | Actual: %s | Expected: %s", actual.Request.PathParams["bar"], expected.Request.PathParams["bar"])
	}
}

func TestOktaClientQueryParam(t *testing.T) {
	url := "https://foo.com"

	actual := BootstrapRequest(url).
		Get("/abc").
		SetQueryParam("foo", "123")

	expected := BootstrapRequest(url)
	expected.Request.URL = url + "/abc/{foo}/{bar}"
	expected.Request.QueryParam = map[string][]string{
		"foo": {"123"},
	}

	if actual.Request.QueryParam.Get("foo") != expected.Request.QueryParam.Get("foo") {
		t.Errorf("FAILED - TestOktaClientQueryParam | Actual: %+v | Expected: %+v", actual.Request.QueryParam.Get("foo"), expected.Request.QueryParam.Get("foo"))
	}
}

func TestOktaClientQueryParamWithMultipleParams(t *testing.T) {
	url := "https://foo.com"

	actual := BootstrapRequest(url).
		Get("/abc").
		SetQueryParam("foo", "123").
		SetQueryParam("bar", "456")

	expected := BootstrapRequest(url)
	expected.Request.URL = url + "/abc/{foo}/{bar}"
	expected.Request.QueryParam = map[string][]string{
		"foo": {"123"},
		"bar": {"456"},
	}

	if actual.Request.QueryParam.Get("foo") != expected.Request.QueryParam.Get("foo") {
		t.Errorf("FAILED - TestOktaClientQueryParamWithMultipleParams | Actual: %+v | Expected: %+v", actual.Request.QueryParam.Get("foo"), expected.Request.QueryParam.Get("foo"))
	}

	if actual.Request.QueryParam.Get("bar") != expected.Request.QueryParam.Get("bar") {
		t.Errorf("FAILED - TestOktaClientQueryParamWithMultipleParams | Actual: %+v | Expected: %+v", actual.Request.QueryParam.Get("bar"), expected.Request.QueryParam.Get("bar"))
	}
}

/********** ROUTE SETTER FUNCTION TESTS **********/

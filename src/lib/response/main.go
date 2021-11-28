package response

import (
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode int
	Body       string
}

func Create() *Response {
	return &Response{}
}

func (response *Response) WithStatus(status int) *Response {
	response.StatusCode = status
	return response
}

func (response *Response) WithBody(body string) *Response {
	response.Body = body
	return response
}

func (response *Response) ToAPIGatewayResponse() (events.APIGatewayProxyResponse, error) {
	defaultHeaders := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Headers:    defaultHeaders,
		Body:       response.Body,
	}, nil
}

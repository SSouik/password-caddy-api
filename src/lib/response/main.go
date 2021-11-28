package response

import (
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode int
	Body       string
}

func Create() Response {
	return Response{}
}

func WithStatus(status int, response *Response) {
	response.StatusCode = status
}

func WithBody(body string, response *Response) {
	response.Body = body
}

func ToAPIGatewayResponse(response Response) events.APIGatewayProxyResponse {
	defaultHeaders := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Headers:    defaultHeaders,
		Body:       response.Body,
	}
}

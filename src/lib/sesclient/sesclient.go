package sesclient

import (
	"context"
	"fmt"
	"password-caddy/api/src/lib/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SesClient struct {
	client *ses.Client
	email  *ses.SendEmailInput
}

type SesResponse struct {
	IsSuccess    bool
	MessageId    string
	ErrorMessage string
}

const OTP_EMAIL_TEMPLATE = `
<h4>Here is your verification token.</h4>
<b>%s</b>
<br/>
<p>
	If you did not recently attempt to login, please ignore this email.
</p>
`

/*
Create a new instance of the AWS Ses Client
*/
func Create(awsConfig aws.Config) *SesClient {
	var client SesClient
	client.client = ses.NewFromConfig(awsConfig)
	return &client
}

/*
Build the email input with the sender and appropriate receiver
*/
func (client *SesClient) BuildEmailRequest(email string) *SesClient {
	var sender string = "me@samuelsouik.com" // update after having password-caddy.com email
	var emails []string = []string{email}
	var charSet string = "UTF-8"
	var subject string = "Verification for Password Caddy"

	code, _ := util.GenerateOTP(6)

	html := fmt.Sprintf(OTP_EMAIL_TEMPLATE, code)

	var input ses.SendEmailInput = ses.SendEmailInput{
		Source: &sender,
		Destination: &types.Destination{
			ToAddresses: emails,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: &charSet,
				Data:    &subject,
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: &charSet,
					Data:    aws.String(html),
				},
			},
		},
	}

	client.email = &input

	return client
}

/*
Send the email via AWS SES
*/
func (client *SesClient) Send() *SesResponse {
	res, err := client.client.SendEmail(context.TODO(), client.email)

	if err != nil {
		return Failure(err.Error())
	}

	return Success(*res.MessageId)
}

/*
Create A successful SES response
*/
func Success(messageId string) *SesResponse {
	return &SesResponse{
		IsSuccess: true,
		MessageId: messageId,
	}
}

/*
Create a failure SES response
*/
func Failure(errorMessage string) *SesResponse {
	return &SesResponse{
		IsSuccess:    false,
		ErrorMessage: errorMessage,
	}
}

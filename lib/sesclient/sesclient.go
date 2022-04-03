package sesclient

import (
	"context"
	"errors"
	"fmt"

	apiTypes "password-caddy/api/core/types"
	"password-caddy/api/lib/util"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SesClient struct {
	Client *ses.Client
	Email  *ses.SendEmailInput
}

type SesResponse struct {
	IsSuccess bool
	Data      interface{}
	Error     apiTypes.PasswordCaddyError
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
	client.Client = ses.NewFromConfig(awsConfig)
	return &client
}

// Verify an email address
func (client *SesClient) SendVerificationEmail(email string) *SesResponse {
	input := ses.VerifyEmailIdentityInput{
		EmailAddress: &email,
	}

	_, err := client.Client.VerifyEmailIdentity(context.TODO(), &input)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(util.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiTypes.PasswordCaddyError{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	return Success("")
}

// Get the verification status of an email address
//
// param: email - string - email address to get verification status of
//
// returns: *SesResponse
func (client *SesClient) GetVerificationStatus(email string) *SesResponse {
	input := ses.GetIdentityVerificationAttributesInput{
		Identities: []string{email},
	}

	response, err := client.Client.GetIdentityVerificationAttributes(context.TODO(), &input)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(util.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiTypes.PasswordCaddyError{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	status := response.VerificationAttributes[email].VerificationStatus

	return Success(struct{ Status types.VerificationStatus }{Status: status})
}

/*
Build the email input with the sender and appropriate receiver
*/
func (client *SesClient) BuildEmailRequest(email, otp string) *SesClient {
	var sender string = "me@samuelsouik.com" // update after having password-caddy.com email
	var emails []string = []string{email}
	var charSet string = "UTF-8"
	var subject string = "Verification for Password Caddy"

	html := fmt.Sprintf(OTP_EMAIL_TEMPLATE, otp)

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

	client.Email = &input

	return client
}

/*
Send the email via AWS SES
*/
func (client *SesClient) Send() *SesResponse {
	res, err := client.Client.SendEmail(context.TODO(), client.Email)

	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) {
			return Failure(util.AWSErrorToPasswordCaddyError(awsErr))
		}

		return Failure(apiTypes.PasswordCaddyError{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	return Success(*res.MessageId)
}

/*
Create A successful SES response
*/
func Success(data interface{}) *SesResponse {
	return &SesResponse{
		IsSuccess: true,
		Data:      data,
	}
}

/*
Create a failure SES response
*/
func Failure(pcError apiTypes.PasswordCaddyError) *SesResponse {
	return &SesResponse{
		IsSuccess: false,
		Error:     pcError,
	}
}

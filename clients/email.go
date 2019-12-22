package clients

import (
	emailService "github.com/obedtandadjaja/project_k_backend/services/email"
)

type EmailClient struct{}

var emailClient *EmailClient

func NewEmailClient() *EmailClient {
	if emailClient != nil {
		return emailClient
	}

	return &EmailClient{}
}

func (emailClient *EmailClient) Send(request *emailService.SendRequest) (*emailService.SendResponse, error) {
	return emailService.NewEmailService().Send(request)
}

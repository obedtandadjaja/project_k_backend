package email

import (
	"context"
	"time"
)

type SendRequest struct {
	Sender        string   `json:"sender"`
	Recipients    []string `json:"recipients"`
	Subject       string   `json:"subject"`
	Body          string   `json:"body"`
	BodyHtml      string   `json:"bodyHtml"`
	BccRecipients []string `json:"bccRecipients"`
}

type SendResponse struct {
	Message string `json:"message"`
	EmailID string `json:"emailId"`
}

func (emailService *EmailService) Send(request *SendRequest) (*SendResponse, error) {
	message := emailService.Mailgun.NewMessage(request.Sender, request.Subject, request.Body, request.Recipients...)
	if request.BodyHtml != "" {
		message.SetHtml(request.BodyHtml)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	msg, id, err := emailService.Mailgun.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	response := SendResponse{
		Message: msg,
		EmailID: id,
	}
	return &response, nil
}

package email

import (
	"github.com/gobuffalo/envy"
	"github.com/mailgun/mailgun-go/v3"
)

type EmailService struct {
	Mailgun *mailgun.MailgunImpl
}

var emailService EmailService

func NewEmailService() *EmailService {
	if emailService.Mailgun == nil {
		mailgunDomain := envy.Get("MAILGUN_DOMAIN", "")
		mailgunApiKey := envy.Get("MAILGUN_API_KEY", "")

		emailService = EmailService{
			Mailgun: mailgun.NewMailgun(mailgunDomain, mailgunApiKey),
		}
	}

	return &emailService
}

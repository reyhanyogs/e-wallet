package service

import (
	"github.com/resend/resend-go/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/internal/config"
)

type emailService struct {
	config *config.Config
}

func NewEmail(config *config.Config) domain.EmailService {
	return &emailService{
		config: config,
	}
}

func (s *emailService) Send(to, subject, body string) error {

	client := resend.NewClient(s.config.Mail.API)

	params := &resend.SendEmailRequest{
		From:    "reyhan@resend.dev",
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	// msg := []byte("" +
	// 	"From: Reyhan Stama <" + s.config.Mail.User + ">\n" +
	// 	"To: " + to + "\n" +
	// 	"Subject: " + subject + "\n" +
	// 	body)

	// auth := smtp.PlainAuth("", s.config.Mail.User, s.config.Mail.Password, s.config.Mail.Host)
	return nil
}

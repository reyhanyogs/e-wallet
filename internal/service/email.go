package service

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

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
	from := mail.Address{"", s.config.Mail.User}
	toMail := mail.Address{"", to}

	// Setup Headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = toMail.String()
	headers["Subject"] = subject

	// Setup Message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := s.config.Mail.Host + ":" + s.config.Mail.Port

	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", s.config.Mail.User, s.config.Mail.Password, host)

	// TLS Config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// From & To
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(toMail.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	return w.Close()
}

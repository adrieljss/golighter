package platform

import (
	"github.com/gofiber/fiber/v3/log"

	"gopkg.in/gomail.v2"
)

type Email struct {
	To      []string
	Subject string
	Body    string
}

type SmtpConfig struct {
	Host         string
	Port         int
	EmailAddress string // alias for Username
	Password     string
}

type SMTPMailer struct {
	config SmtpConfig
}

// inits SMTPMailer with config
// also validates connection to SMTP server
func NewMailer(config SmtpConfig) *SMTPMailer {
	mailer := &SMTPMailer{
		config: config,
	}

	mailer.ValidateConnection()

	return mailer
}

func (m *SMTPMailer) Send(email *Email) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.EmailAddress)
	msg.SetHeader("To", email.To...)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/html", email.Body)

	dialer := gomail.NewDialer(m.config.Host, m.config.Port, m.config.EmailAddress, m.config.Password)

	if err := dialer.DialAndSend(msg); err != nil {
		log.Errorf("failed to send email: %s", err.Error())
		return err
	}

	return nil
}

// validate connection to SMTP server on startup
// panics if connection fails
func (m *SMTPMailer) ValidateConnection() {
	dialer := gomail.NewDialer(m.config.Host, m.config.Port, m.config.EmailAddress, m.config.Password)

	sender, err := dialer.Dial()
	if err != nil {
		log.Fatalf("failed to connect to SMTP server: %s", err.Error())
	}
	defer sender.Close()

	log.Info("successfully connected to SMTP server")
}

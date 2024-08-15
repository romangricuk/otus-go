package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
)

type Client interface {
	SendEmail(to string, subject string, body string) error
}

type smtpClient struct {
	config *config.EmailConfig
}

func NewSMTPClient(config *config.EmailConfig) Client {
	fmt.Printf("config: %+v\n", config)
	return &smtpClient{
		config: config,
	}
}

func (c *smtpClient) SendEmail(to string, subject string, body string) error {
	auth := smtp.CRAMMD5Auth(c.config.Username, c.config.Password)
	msg := []byte("From: " + c.config.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", c.config.SMTPServer, c.config.SMTPPort)

	if c.config.UseTLS { //nolint:nestif
		// Используем TLS-соединение
		tlsconfig := &tls.Config{
			InsecureSkipVerify: c.config.InsecureSkipVerify, //nolint:gosec
			ServerName:         c.config.SMTPServer,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to establish TLS connection: %w", err)
		}

		client, err := smtp.NewClient(conn, c.config.SMTPServer)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer client.Quit()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err = client.Mail(c.config.From); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to add recipient: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to start data command: %w", err)
		}

		if _, err = w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		if err = w.Close(); err != nil {
			return fmt.Errorf("failed to close message writer: %w", err)
		}

		return client.Quit()
	} else { //nolint:revive
		// Используем обычное нешифрованное соединение
		if err := smtp.SendMail(addr, auth, c.config.From, []string{to}, msg); err != nil {
			return fmt.Errorf("on send email: %w", err)
		}
	}

	return nil
}

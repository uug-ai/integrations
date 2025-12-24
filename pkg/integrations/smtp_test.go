package integrations

import (
	"os"
	"testing"
)

func setupSMTPTest() SMTP {
	// Initialize SMTP configuration from environment variables
	smtpMailtrap := SMTP{
		Server:    os.Getenv("SMTP_SERVER"),
		Port:      os.Getenv("SMTP_PORT"),
		Username:  os.Getenv("SMTP_USERNAME"),
		Password:  os.Getenv("SMTP_PASSWORD"),
		EmailFrom: os.Getenv("EMAIL_FROM"),
		EmailTo:   os.Getenv("EMAIL_TO"),
	}
	return smtpMailtrap
}

func TestValidSMTP(t *testing.T) {
	smtpMailtrap := setupSMTPTest()

	// Send message to SMTP server.
	err := smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

func TestWrongEmail(t *testing.T) {
	smtpMailtrap := setupSMTPTest()
	smtpMailtrap.EmailTo = "invalid-email-address"

	// Send message to SMTP server.
	err := smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")
	if err == nil {
		t.Errorf("expected error got nil")
	}
}

func TestWrongServer(t *testing.T) {
	smtpMailtrap := setupSMTPTest()
	smtpMailtrap.Server = "wrong.smtp.server"

	// Send message to SMTP server.
	err := smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")
	if err == nil {
		t.Errorf("expected error got nil")
	}
}

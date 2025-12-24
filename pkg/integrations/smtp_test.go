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

func TestSMTPValidation(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*SMTP)
		expectError bool
	}{
		{
			name: "MissingProperties",
			setup: func(s *SMTP) {
				s.EmailFrom = ""
				s.EmailTo = "cedric@lol.be"
			},
			expectError: true,
		},
		{
			name: "WrongEmail",
			setup: func(s *SMTP) {
				s.EmailTo = "invalid-email-address"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smtpMailtrap := setupSMTPTest()
			tt.setup(&smtpMailtrap)

			// Send message to SMTP server.
			err := smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")

			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

func TestSMTPServer(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*SMTP)
		expectError bool
	}{
		{
			name:        "ValidSMTP",
			setup:       func(s *SMTP) {},
			expectError: false,
		},
		{
			name: "WrongServer",
			setup: func(s *SMTP) {
				s.Server = "wrong.smtp.server"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smtpMailtrap := setupSMTPTest()
			tt.setup(&smtpMailtrap)

			// Send message to SMTP server.
			err := smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")

			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

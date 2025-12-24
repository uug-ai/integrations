package integrations

import (
	"os"
	"strconv"
	"testing"
)

func setupSMTPTest() (*SMTP, error) {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpMailtrap, err := CreateSMTP(
		WithSMTPServer(os.Getenv("SMTP_SERVER")),
		WithSMTPPort(port),
		WithSMTPUsername(os.Getenv("SMTP_USERNAME")),
		WithSMTPPassword(os.Getenv("SMTP_PASSWORD")),
		WithSMTPEmailFrom(os.Getenv("EMAIL_FROM")),
		WithSMTPEmailTo(os.Getenv("EMAIL_TO")),
	)
	return smtpMailtrap, err
}

func TestSMTPValidation(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*SMTP)
		expectError bool
	}{
		{
			name: "MissingEmailFrom",
			setup: func(s *SMTP) {
				s.EmailFrom = ""
			},
			expectError: true,
		},
		{
			name: "MissingEmailTo",
			setup: func(s *SMTP) {
				s.EmailTo = ""
			},
			expectError: true,
		},
		{
			name: "WrongEmailTo",
			setup: func(s *SMTP) {
				s.EmailTo = "invalid-email-address"
			},
			expectError: true,
		},
		{
			name: "WrongEmailFrom",
			setup: func(s *SMTP) {
				s.EmailFrom = "not-an-email"
			},
			expectError: true,
		},
		{
			name: "MissingServer",
			setup: func(s *SMTP) {
				s.Server = ""
			},
			expectError: true,
		},
		{
			name: "MissingPort",
			setup: func(s *SMTP) {
				s.Port = 0
			},
			expectError: true,
		},
		{
			name: "NegativePort",
			setup: func(s *SMTP) {
				s.Port = -25
			},
			expectError: true,
		},
		{
			name: "ValidWithoutCredentials",
			setup: func(s *SMTP) {
				s.Username = ""
				s.Password = ""
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smtpMailtrap, err := setupSMTPTest()
			if err != nil {
				t.Fatalf("failed to setup SMTP: %v", err)
			}
			tt.setup(smtpMailtrap)

			err = smtpMailtrap.Validate()
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			// Send message to SMTP server.
			err = smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")

			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

func TestSMTPFieldEmpty(t *testing.T) {

	smtpMailtrap, err := setupSMTPTest()
	if err != nil {
		t.Fatalf("failed to setup SMTP: %v", err)
	}

	tests := []struct {
		title       string
		body        string
		textBody    string
		expectError bool
	}{
		{"", "Body", "<p>Body</p>", true},
		{"Title", "", "<p>Body</p>", true},
		{"Title", "Body", "", true},
		{"Title", "Body", "<p>Body</p>", false},
	}

	for _, tt := range tests {
		err := smtpMailtrap.Send(tt.title, tt.body, tt.textBody)
		if tt.expectError && err == nil {
			t.Errorf("expected error got nil for title: '%s', body: '%s', textBody: '%s'", tt.title, tt.body, tt.textBody)
		}
		if !tt.expectError && err != nil {
			t.Errorf("expected no error got %v for title: '%s', body: '%s', textBody: '%s'", err, tt.title, tt.body, tt.textBody)
		}
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
		{
			name: "WrongPort",
			setup: func(s *SMTP) {
				s.Port = -100
			},
			expectError: true,
		},
		{
			name: "InvalidPort",
			setup: func(s *SMTP) {
				s.Port = 0
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smtpMailtrap, err := setupSMTPTest()
			if err != nil {
				t.Fatalf("failed to setup SMTP: %v", err)
			}
			tt.setup(smtpMailtrap)

			// Send message to SMTP server.
			err = smtpMailtrap.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

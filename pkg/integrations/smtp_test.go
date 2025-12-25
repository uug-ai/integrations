package integrations

import (
	"os"
	"strconv"
	"testing"
)

func setupSMTPTest() (*SMTP, error) {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	opts := NewSMTPOptions().
		Server(os.Getenv("SMTP_SERVER")).
		Port(port).
		Username(os.Getenv("SMTP_USERNAME")).
		Password(os.Getenv("SMTP_PASSWORD")).
		From(os.Getenv("EMAIL_FROM")).
		To(os.Getenv("EMAIL_TO")).
		Build()

	smtp, err := NewSMTP(opts)
	return smtp, err
}

func TestSMTPValidation(t *testing.T) {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	tests := []struct {
		name        string
		buildOpts   func() *SMTPOptions
		expectError bool
	}{
		{
			name: "MissingEmailFrom",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingEmailTo",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongEmailTo",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To("invalid-email-address").
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongEmailFrom",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From("not-an-email").
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingServer",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "NegativePort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(-25).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "ValidWithoutCredentials",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()

			// Try to create SMTP client - validation happens here
			_, err := NewSMTP(opts)

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
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	tests := []struct {
		name        string
		buildOpts   func() *SMTPOptions
		expectError bool
	}{
		{
			name: "ValidSMTP",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: false,
		},
		{
			name: "WrongServer",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server("wrong.smtp.server").
					Port(port).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(-100).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					Server(os.Getenv("SMTP_SERVER")).
					Port(0).
					Username(os.Getenv("SMTP_USERNAME")).
					Password(os.Getenv("SMTP_PASSWORD")).
					From(os.Getenv("EMAIL_FROM")).
					To(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()

			// Try to create SMTP client
			smtpClient, err := NewSMTP(opts)

			// For validation errors, check client creation
			if tt.expectError && tt.name != "ValidSMTP" && tt.name != "WrongServer" {
				if err == nil {
					t.Errorf("expected error during client creation got nil")
				}
				return
			}

			if err != nil && !tt.expectError {
				t.Fatalf("failed to create SMTP client: %v", err)
			}

			// For runtime errors (like wrong server), try to send
			if smtpClient != nil {
				err = smtpClient.Send("Test Subject", "This is the body of the email.", "<p>This is the body of the email.</p>")
				if tt.expectError && err == nil {
					t.Errorf("expected error got nil")
				}
				if !tt.expectError && err != nil {
					t.Errorf("expected error to be nil got %v", err)
				}
			}
		})
	}
}

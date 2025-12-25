package integrations

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"gopkg.in/gomail.v2"
)

// MockMailDialer is a mock implementation of MailDialer for testing
type MockMailDialer struct {
	DialAndSendFunc func(m ...*gomail.Message) error
	DialFunc        func() (gomail.SendCloser, error)
	DialCalled      bool
	SendCalled      bool
}

func (m *MockMailDialer) DialAndSend(msgs ...*gomail.Message) error {
	m.SendCalled = true
	if m.DialAndSendFunc != nil {
		return m.DialAndSendFunc(msgs...)
	}
	return nil
}

func (m *MockMailDialer) Dial() (gomail.SendCloser, error) {
	m.DialCalled = true
	if m.DialFunc != nil {
		return m.DialFunc()
	}
	return nil, nil
}

func setupSMTPTest() (*SMTP, error) {
	opts := NewSMTPOptions().
		SetServer(os.Getenv("SMTP_SERVER")).
		SetPort(func() int { port, _ := strconv.Atoi(os.Getenv("SMTP_PORT")); return port }()).
		SetUsername(os.Getenv("SMTP_USERNAME")).
		SetPassword(os.Getenv("SMTP_PASSWORD")).
		SetFrom(os.Getenv("EMAIL_FROM")).
		SetTo(os.Getenv("EMAIL_TO")).
		Build()

	smtp, err := NewSMTP(opts, nil) // nil uses default production dialer
	return smtp, err
}

func TestSMTPValidation(t *testing.T) {
	// Use mock dialer to avoid actual network calls
	mockDialer := &MockMailDialer{}

	tests := []struct {
		name        string
		buildOpts   func() *SMTPOptions
		expectError bool
	}{
		{
			name: "MissingEmailFrom",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingEmailTo",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongEmailTo",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo("invalid-email-address").
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongEmailFrom",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom("not-an-email").
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingServer",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "NegativePort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(-25).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "ValidWithoutCredentials",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()

			// Try to create SMTP client - validation happens here
			_, err := NewSMTP(opts, mockDialer)

			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}
func TestSMTPSendWithMock(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		body        string
		textBody    string
		dialError   error
		sendError   error
		expectError bool
		expectDial  bool
		expectSend  bool
	}{
		{
			name:        "EmptyTitle",
			title:       "",
			body:        "Body",
			textBody:    "<p>Body</p>",
			expectError: true,
			expectDial:  false,
			expectSend:  false,
		},
		{
			name:        "EmptyBody",
			title:       "Title",
			body:        "",
			textBody:    "<p>Body</p>",
			expectError: true,
			expectDial:  false,
			expectSend:  false,
		},
		{
			name:        "EmptyTextBody",
			title:       "Title",
			body:        "Body",
			textBody:    "",
			expectError: true,
			expectDial:  false,
			expectSend:  false,
		},
		{
			name:        "DialError",
			title:       "Title",
			body:        "Body",
			textBody:    "<p>Body</p>",
			dialError:   errors.New("dial failed"),
			expectError: true,
			expectDial:  true,
			expectSend:  false,
		},
		{
			name:        "SendError",
			title:       "Title",
			body:        "Body",
			textBody:    "<p>Body</p>",
			sendError:   errors.New("send failed"),
			expectError: true,
			expectDial:  true,
			expectSend:  true,
		},
		{
			name:        "Success",
			title:       "Title",
			body:        "Body",
			textBody:    "<p>Body</p>",
			expectError: false,
			expectDial:  true,
			expectSend:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock dialer
			mockDialer := &MockMailDialer{
				DialFunc: func() (gomail.SendCloser, error) {
					return nil, tt.dialError
				},
				DialAndSendFunc: func(m ...*gomail.Message) error {
					return tt.sendError
				},
			}

			// Create SMTP client with mock
			opts := NewSMTPOptions().
				SetServer("smtp.test.com").
				SetPort(587).
				SetUsername("user").
				SetPassword("pass").
				SetFrom("from@test.com").
				SetTo("to@test.com").
				Build()

			smtp, err := NewSMTP(opts, mockDialer)
			if err != nil {
				t.Fatalf("failed to create SMTP client: %v", err)
			}

			// Test Send
			err = smtp.Send(tt.title, tt.body, tt.textBody)

			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error got %v", err)
			}

			if mockDialer.DialCalled != tt.expectDial {
				t.Errorf("expected Dial called=%v, got=%v", tt.expectDial, mockDialer.DialCalled)
			}

			if mockDialer.SendCalled != tt.expectSend {
				t.Errorf("expected Send called=%v, got=%v", tt.expectSend, mockDialer.SendCalled)
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
		buildOpts   func() *SMTPOptions
		expectError bool
	}{
		{
			name: "ValidSMTP",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: false,
		},
		{
			name: "WrongServer",
			buildOpts: func() *SMTPOptions {
				port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
				return NewSMTPOptions().
					SetServer("wrong.smtp.server").
					SetPort(port).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "WrongPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(-100).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidPort",
			buildOpts: func() *SMTPOptions {
				return NewSMTPOptions().
					SetServer(os.Getenv("SMTP_SERVER")).
					SetPort(0).
					SetUsername(os.Getenv("SMTP_USERNAME")).
					SetPassword(os.Getenv("SMTP_PASSWORD")).
					SetFrom(os.Getenv("EMAIL_FROM")).
					SetTo(os.Getenv("EMAIL_TO")).
					Build()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()

			// Try to create SMTP client
			smtpClient, err := NewSMTP(opts, nil) // nil uses default production dialer

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

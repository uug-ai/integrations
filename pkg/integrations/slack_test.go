package integrations

import (
	"errors"
	"os"
	"testing"

	"github.com/slack-go/slack"
)

// MockSlackWebhookClient is a mock implementation of SlackWebhookClient for testing
type MockSlackWebhookClient struct {
	PostWebhookFunc func(url string, msg *slack.WebhookMessage) error
	PostCalled      bool
	LastURL         string
	LastMessage     *slack.WebhookMessage
}

func (m *MockSlackWebhookClient) PostWebhook(url string, msg *slack.WebhookMessage) error {
	m.PostCalled = true
	m.LastURL = url
	m.LastMessage = msg
	if m.PostWebhookFunc != nil {
		return m.PostWebhookFunc(url, msg)
	}
	return nil
}

func setupSlackTest() (*Slack, error) {
	opts := NewSlackOptions().
		SetHook(os.Getenv("SLACK_HOOK")).
		SetUsername(os.Getenv("SLACK_USERNAME")).
		Build()

	slack, err := NewSlack(opts, nil) // nil uses default production client
	return slack, err
}

func TestSlackValidation(t *testing.T) {
	// Use mock client to avoid actual network calls
	mockClient := &MockSlackWebhookClient{}

	tests := []struct {
		name        string
		buildOpts   func() *SlackOptions
		expectError bool
	}{
		{
			name: "MissingHook",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: true,
		},
		{
			name: "MissingUsername",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook(os.Getenv("SLACK_HOOK")).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidHookURL",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook("not-a-valid-url").
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: true,
		},
		{
			name: "ValidSlack",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook(os.Getenv("SLACK_HOOK")).
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()
			_, err := NewSlack(opts, mockClient)
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

func TestSlackFieldEmpty(t *testing.T) {
	// Use mock client to avoid actual network calls
	mockClient := &MockSlackWebhookClient{}

	opts := NewSlackOptions().
		SetHook(os.Getenv("SLACK_HOOK")).
		SetUsername(os.Getenv("SLACK_USERNAME")).
		Build()

	slackIntegration, err := NewSlack(opts, mockClient)
	if err != nil {
		t.Fatalf("failed to setup Slack: %v", err)
	}

	tests := []struct {
		body        string
		url         string
		expectError bool
	}{
		{body: "", url: "https://example.com", expectError: true},
		{body: "Test message", url: "", expectError: false},
		{body: "Test message", url: "https://example.com", expectError: false},
	}

	for _, tt := range tests {
		mockClient.PostCalled = false // Reset for each test
		err := slackIntegration.Send(tt.body, tt.url)
		if tt.expectError && err == nil {
			t.Errorf("expected error got nil for body: '%s', url: '%s'", tt.body, tt.url)
		}
		if !tt.expectError && err != nil {
			t.Errorf("expected no error got %v for body: '%s', url: '%s'", err, tt.body, tt.url)
		}
		if !tt.expectError && !mockClient.PostCalled {
			t.Errorf("expected PostWebhook to be called but it wasn't")
		}
	}
}

func TestSlackChannel(t *testing.T) {
	tests := []struct {
		name          string
		mockPostError error
		expectError   bool
	}{
		{
			name:          "ValidSlack",
			mockPostError: nil,
			expectError:   false,
		},
		{
			name:          "WrongHook",
			mockPostError: errors.New("invalid webhook URL"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client with specific error behavior
			mockClient := &MockSlackWebhookClient{
				PostWebhookFunc: func(url string, msg *slack.WebhookMessage) error {
					return tt.mockPostError
				},
			}

			opts := NewSlackOptions().
				SetHook(os.Getenv("SLACK_HOOK")).
				SetUsername(os.Getenv("SLACK_USERNAME")).
				Build()

			slackIntegration, err := NewSlack(opts, mockClient)
			if err != nil {
				t.Fatalf("failed to setup Slack: %v", err)
			}

			// Send message to Slack channel.
			err = slackIntegration.Send("Test message from UUG AI", "https://uug.ai")
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if !mockClient.PostCalled {
				t.Errorf("expected PostWebhook to be called but it wasn't")
			}
		})
	}
}

func TestIntegrationSlackWebhook(t *testing.T) {
	tests := []struct {
		name        string
		buildOpts   func() *SlackOptions
		expectError bool
	}{
		{
			name: "ValidSlack",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook(os.Getenv("SLACK_HOOK")).
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: false,
		},
		{
			name: "WrongHook",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook("https://hooks.slack.com/services/WRONG/HOOK/URL").
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidHookFormat",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook("not-a-valid-url").
					SetUsername(os.Getenv("SLACK_USERNAME")).
					Build()
			},
			expectError: true,
		},
		{
			name: "EmptyUsername",
			buildOpts: func() *SlackOptions {
				return NewSlackOptions().
					SetHook(os.Getenv("SLACK_HOOK")).
					Build()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()

			// Try to create Slack client
			slackClient, err := NewSlack(opts) // no client uses default production client

			// For validation errors, check client creation
			if tt.expectError && (tt.name == "InvalidHookFormat" || tt.name == "EmptyUsername") {
				if err == nil {
					t.Errorf("expected error during client creation got nil")
				}
				return
			}

			if err != nil && !tt.expectError {
				t.Fatalf("failed to create Slack client: %v", err)
			}

			// For runtime errors (like wrong hook), try to send
			if slackClient != nil {
				err = slackClient.Send("Test message from integration test", "https://example.com")
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

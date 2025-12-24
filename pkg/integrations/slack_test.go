package integrations

import (
	"os"
	"testing"
)

func setupSlackTest() (*Slack, error) {
	slackIntegration, err := CreateSlack(
		WithSlackHook(os.Getenv("SLACK_HOOK")),
		WithSlackUsername(os.Getenv("SLACK_USERNAME")),
	)
	return slackIntegration, err
}

func TestSlackValidation(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*Slack)
		expectError bool
	}{
		{
			name: "MissingHook",
			setup: func(s *Slack) {
				s.Hook = ""
			},
			expectError: true,
		},
		{
			name: "MissingUsername",
			setup: func(s *Slack) {
				s.Username = ""
			},
			expectError: true,
		},
		{
			name: "InvalidHookURL",
			setup: func(s *Slack) {
				s.Hook = "not-a-valid-url"
			},
			expectError: true,
		},
		{
			name: "ValidSlack",
			setup: func(s *Slack) {
				// No changes - should be valid
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slackIntegration, err := setupSlackTest()
			if err != nil {
				t.Fatalf("failed to setup Slack: %v", err)
			}
			tt.setup(slackIntegration)

			err = slackIntegration.Validate()
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
	slackIntegration, err := setupSlackTest()
	if err != nil {
		t.Fatalf("failed to setup Slack: %v", err)
	}

	tests := []struct {
		body        string
		url         string
		expectError bool
	}{
		{"", "https://example.com", true},
		{"Test message", "", false},
		{"Test message", "https://example.com", false},
	}

	for _, tt := range tests {
		err := slackIntegration.Send(tt.body, tt.url)
		if tt.expectError && err == nil {
			t.Errorf("expected error got nil for body: '%s', url: '%s'", tt.body, tt.url)
		}
		if !tt.expectError && err != nil {
			t.Errorf("expected no error got %v for body: '%s', url: '%s'", err, tt.body, tt.url)
		}
	}
}

func TestSlackChannel(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*Slack)
		expectError bool
	}{
		{
			name:        "ValidSlack",
			setup:       func(s *Slack) {},
			expectError: false,
		},
		{
			name: "WrongHook",
			setup: func(s *Slack) {
				s.Hook = "https://hooks.slack.com/services/WRONG/HOOK/URL"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slackIntegration, err := setupSlackTest()
			if err != nil {
				t.Fatalf("failed to setup Slack: %v", err)
			}
			tt.setup(slackIntegration)

			// Send message to Slack channel.
			err = slackIntegration.Send("Test message from UUG AI", "https://uug.ai")
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

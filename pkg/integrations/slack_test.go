package integrations

import (
	"os"
	"testing"
)

func setupSlackTest() Slack {
	// Initialize Slack integration from environment variables
	slackIntegration := Slack{
		Hook:     os.Getenv("SLACK_HOOK"),
		Username: os.Getenv("SLACK_USERNAME"),
	}
	return slackIntegration
}

func TestSlackChannel(t *testing.T) {
	slackIntegration := setupSlackTest()
	// Send message to Slack channel.
	err := slackIntegration.Send("Test message from UUG AI", "https://uug.ai")
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

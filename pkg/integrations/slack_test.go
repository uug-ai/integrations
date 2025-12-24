package integrations

import (
	"os"
	"testing"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

var slackIntegration Slack
var slackTimeout = 0
var slackTimeoutIncrement = 1000

func setupSlackTest() Slack {
	// Initialize Slack integration from environment variables
	slackIntegration = Slack{
		Hook:     os.Getenv("SLACK_HOOK"),
		Username: os.Getenv("SLACK_USERNAME"),
	}

	// Timeout, to avoid hitting issues with Slack API.
	slackTimeout = slackTimeout + slackTimeoutIncrement
	tout := time.Duration(slackTimeout) * time.Millisecond
	time.Sleep(tout)

	return slackIntegration
}

func TestSlackChannel(t *testing.T) {
	slackIntegration := setupSlackTest()

	// Message to send to the Slack channel.
	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "Alert: Kerberos Hub detected something."
	m.Body = "Alert: Kerberos Hub detected something."
	m.User = "cedricve"
	m.UserId = "23235235235235"
	m.SequenceId = "5a72d0f6e17699d18adb5e17"
	m.Unread = true
	m.Media = []models.Media{}
	m.Media = append(m.Media, models.Media{
		StartTimestamp: 1670618365,
		AtRuntimeMetadata: &models.MediaAtRuntimeMetadata{
			VideoUrl:     "https://example.com/video.mp4",
			ThumbnailUrl: "https://example.com/thumbnail.jpg",
		},
	})

	// Send message to Slack channel.
	err := slackIntegration.Send(m)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

package test

import (
	"testing"
	"time"

	channels "github.com/uug-ai/hub-pipeline-notification/channels"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

func TestSlackChannel(t *testing.T) {
	m := message.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "Alert: Kerberos Hub detected something."
	m.Body = "Alert: Kerberos Hub detected something."
	m.User = "cedricve"
	m.UserId = "23235235235235"
	m.SequenceId = "5a72d0f6e17699d18adb5e17"
	m.Unread = true
	m.Media = []message.Media{}
	m.Media = append(m.Media, message.Media{
		Timestamp:    1670618365,
		Type:         "video",
		Url:          "",
		ThumbnailUrl: "",
	})

	// Get User notification channels.
	// ....

	// Send message to all channels.
	slack := channels.Slack{
		Hook:     "https://hooks.slack.com/services/xxxx/xxx/xxxx",
		Username: "UUG.AI Bot",
	}
	err := slack.Send(m)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

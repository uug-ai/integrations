package integrations

import (
	"testing"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

func TestSlackChannel(t *testing.T) {
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
	})

	// Get User notification channels.
	// ....

	// Send message to all channels.
	slack := Slack{
		Hook:     "https://hooks.slack.com/services/xxxx/xxx/xxxx",
		Username: "UUG.AI Bot",
	}
	err := slack.Send(m)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

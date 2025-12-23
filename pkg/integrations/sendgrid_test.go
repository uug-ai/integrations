package integrations

import (
	"testing"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

func TestSendgrid(t *testing.T) {
	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "Your user has been disabled."
	m.Body = "zgeezgzegezzgeezzge"
	m.User = "cedricve"
	m.Email = "xxx@xxxx.io"
	m.UserId = "23235235235235"
	m.SequenceId = "5a72d0f6e17699d18adb5e17"
	m.DataUsage = "124124"
	m.Unread = true

	// Send message to all channels.
	sendgrid := Sendgrid{
		TemplateId: "c4ed7742-8742-4f6d-bfef-b3e42a62ebca",
		EmailFrom:  "xxx@xxx.io",
		EmailTo:    "xxx@xxx.io",
	}
	sendgrid.SendNotification(m)
}

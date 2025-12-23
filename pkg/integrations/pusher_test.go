package integrations

import (
	"testing"
	"time"
)

func TestPusher(t *testing.T) {
	m := Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "Something happened"
	m.Body = "zgeezgzegezzgeezzge"
	m.User = "cedricve"
	m.UserId = "23235235235235"
	m.SequenceId = "5a72d0f6e17699d18adb5e17"
	m.Unread = true

	// Get User notification channels.
	// ....

	// Send message to all channels.
	pusher := Pusher{
		Channel: "notification",
	}
	pusher.Send(m)
}

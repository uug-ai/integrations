package test

import (
	"time"

	channels "github.com/uug-ai/hub-pipeline-notification/channels"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

func main() {
	m := message.Message{}
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
	pusher := channels.Pusher{
		Channel: "notification",
	}
	pusher.Send(m)

	return
}

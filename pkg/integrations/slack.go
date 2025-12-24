package integrations

import (
	"github.com/slack-go/slack"
)

type Slack struct {
	Token    string `json:"token,omitempty"`   // Old
	Channel  string `json:"channel,omitempty"` // Old
	Hook     string `json:"hook,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s Slack) Send(body string, url string) error {

	hook := s.Hook
	username := s.Username

	text := body
	if url != "" {
		text = text + "\r\n" + url
	}

	attachment := slack.Attachment{
		Color:    "good",
		ImageURL: url,
	}

	msg := slack.WebhookMessage{
		Username:    username,
		Text:        text,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(hook, &msg)
	return err
}

package integrations

import (
	"github.com/slack-go/slack"
	"github.com/uug-ai/models/pkg/models"
)

type Slack struct {
	Token    string `json:"token,omitempty"`   // Old
	Channel  string `json:"channel,omitempty"` // Old
	Hook     string `json:"hook,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s Slack) Send(message models.Message) error {

	hook := s.Hook
	username := s.Username

	url := ""
	if len(message.Media) > 0 {
		longUrl := message.Media[0].AtRuntimeMetadata.VideoUrl
		url = longUrl
	}

	text := message.Body
	if url != "" {
		text = text + "\r\n" + url
	}

	attachment := slack.Attachment{
		Color:    "good",
		ImageURL: message.Media[0].AtRuntimeMetadata.ThumbnailUrl,
	}

	msg := slack.WebhookMessage{
		Username:    username,
		Text:        text,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(hook, &msg)
	return err
}

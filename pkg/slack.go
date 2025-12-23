package pkg

import (
	"github.com/slack-go/slack"
)

type Slack struct {
	Token    string `json:"token,omitempty"`   // Old
	Channel  string `json:"channel,omitempty"` // Old
	Hook     string `json:"hook,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s Slack) Send(message Message) error {

	//token := s.Token
	//channelName := s.Channel
	hook := s.Hook
	username := s.Username

	url := ""
	if len(message.Media) > 0 {
		longUrl := message.Media[0].Url
		url = longUrl
		/*provider := "tinyurl"
		shortenedUrl, err := shorturl.Shorten(longUrl, provider)
		if err == nil {
			url = string(shortenedUrl)
		}*/
	}

	text := message.Body
	if url != "" {
		text = text + "\r\n" + url
	}

	attachment := slack.Attachment{
		Color:    "good",
		ImageURL: message.Media[0].ThumbnailUrl,
	}

	msg := slack.WebhookMessage{
		Username:    username,
		Text:        text,
		Attachments: []slack.Attachment{attachment},
	}
	err := slack.PostWebhook(hook, &msg)
	return err
}

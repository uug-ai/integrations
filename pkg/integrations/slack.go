package integrations

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/slack-go/slack"
)

type Slack struct {
	Hook     string `json:"hook" validate:"required,url"`
	Username string `json:"username " validate:"required"`
}

// WithSlackHook sets the Slack webhook URL
func WithSlackHook(hook string) Option[Slack] {
	return func(s *Slack) {
		s.Hook = hook
	}
}

// WithSlackUsername sets the Slack username
func WithSlackUsername(username string) Option[Slack] {
	return func(s *Slack) {
		s.Username = username
	}
}

// CreateSlack creates a new Slack instance with the provided options
func CreateSlack(opts ...Option[Slack]) (*Slack, error) {
	slack := &Slack{}

	// Apply all options
	for _, opt := range opts {
		opt(slack)
	}

	// Validate SMTP configuration
	err := slack.Validate()
	if err != nil {
		return nil, err
	}

	return slack, nil
}

func (slack *Slack) Validate() error {
	validate := validator.New()
	err := validate.Struct(slack)
	if err != nil {
		return err
	}
	return nil
}

func (s *Slack) Send(body string, url string) error {
	hook := s.Hook
	username := s.Username

	if body == "" {
		return errors.New("message body is empty")
	}

	// Append URL to the message body if provided
	text := body
	if url != "" {
		text = text + "\r\n" + url
	}

	// Create Slack webhook message
	msg := slack.WebhookMessage{
		Username: username,
		Text:     text,
		Attachments: []slack.Attachment{
			{
				Color:    "good",
				ImageURL: url,
			},
		},
	}
	err := slack.PostWebhook(hook, &msg)
	return err
}

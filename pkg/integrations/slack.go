package integrations

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/slack-go/slack"
)

// SlackWebhookClient is an interface for posting messages to Slack
type SlackWebhookClient interface {
	PostWebhook(url string, msg *slack.WebhookMessage) error
}

// SlackWebhookClientImpl is the default implementation using slack library
type SlackWebhookClientImpl struct{}

// PostWebhook implements SlackWebhookClient interface
func (s *SlackWebhookClientImpl) PostWebhook(url string, msg *slack.WebhookMessage) error {
	return slack.PostWebhook(url, msg)
}

// NewSlackWebhookClient creates a new default Slack webhook client
func NewSlackWebhookClient() SlackWebhookClient {
	return &SlackWebhookClientImpl{}
}

// SlackOptions holds the configuration for Slack
type SlackOptions struct {
	Hook     string `validate:"required,url"`
	Username string `validate:"required"`
}

// SlackOptionsBuilder provides a fluent interface for building Slack options
type SlackOptionsBuilder struct {
	options *SlackOptions
}

// NewSlackOptions creates a new Slack options builder
func NewSlackOptions() *SlackOptionsBuilder {
	return &SlackOptionsBuilder{
		options: &SlackOptions{},
	}
}

// SetHook sets the Slack webhook URL
func (b *SlackOptionsBuilder) SetHook(hook string) *SlackOptionsBuilder {
	b.options.Hook = hook
	return b
}

// SetUsername sets the Slack username
func (b *SlackOptionsBuilder) SetUsername(username string) *SlackOptionsBuilder {
	b.options.Username = username
	return b
}

// Build returns the configured SlackOptions
func (b *SlackOptionsBuilder) Build() *SlackOptions {
	return b.options
}

// Slack represents a Slack client instance
type Slack struct {
	options *SlackOptions
	client  SlackWebhookClient
}

// NewSlack creates a new Slack client with the provided options
// If client is not provided, a default SlackWebhookClient will be created
func NewSlack(opts *SlackOptions, client ...SlackWebhookClient) (*Slack, error) {
	// Validate Slack configuration
	validate := validator.New()
	err := validate.Struct(opts)
	if err != nil {
		return nil, err
	}

	// If no client provided, create default production client
	var c SlackWebhookClient
	if len(client) == 0 {
		c = NewSlackWebhookClient()
	} else {
		c = client[0]
	}

	return &Slack{
		options: opts,
		client:  c,
	}, nil
}

// Send sends a message to Slack using the configured webhook
// Parameters:
//   - body: The message text to send
//   - url: An optional URL to append to the message and include as an image attachment
//
// Returns:
//   - error: An error if body is empty or if posting to Slack fails
func (s *Slack) Send(body string, url string) error {
	if body == "" {
		return errors.New("message body is empty")
	}

	// Append URL to the message body if provided
	text := body
	if url != "" {
		text = text + "\r\n" + url
	}

	// Create Slack webhook message
	msg := &slack.WebhookMessage{
		Username: s.options.Username,
		Text:     text,
		Attachments: []slack.Attachment{
			{
				Color:    "good",
				ImageURL: url,
			},
		},
	}

	return s.client.PostWebhook(s.options.Hook, msg)
}

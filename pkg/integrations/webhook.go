package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/uug-ai/integrations/pkg/message"
)

// WebhookOptions configures Webhook.
type WebhookOptions struct {
	// Url is the URL the message is posted to.
	Url string `json:"url,omitempty" validate:"required,url"`
}

// WebhookOptionsBuilder builds WebhookOptions.
type WebhookOptionsBuilder struct {
	options *WebhookOptions
}

// NewWebhookOptions starts building WebhookOptions.
func NewWebhookOptions() *WebhookOptionsBuilder {
	return &WebhookOptionsBuilder{options: &WebhookOptions{}}
}

// SetUrl sets the URL the message is posted to.
func (b *WebhookOptionsBuilder) SetUrl(url string) *WebhookOptionsBuilder {
	b.options.Url = url
	return b
}

// Build returns the configured WebhookOptions.
func (b *WebhookOptionsBuilder) Build() *WebhookOptions {
	return b.options
}

// Webhook posts the message as JSON to a URL.
type Webhook struct {
	options *WebhookOptions
}

// NewWebhook validates opts and returns a Webhook sender that uses a copy of them.
func NewWebhook(opts *WebhookOptions) (*Webhook, error) {
	if err := validateOptions("webhook", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Webhook{options: &options}, nil
}

func (webhook *Webhook) Send(message message.Message) bool {
	if err := webhook.SendWithContext(context.Background(), message); err != nil {
		log.Print(err)
		return false
	}
	return true
}

type webhookResponseError struct {
	statusCode int
}

func (err webhookResponseError) Error() string {
	return fmt.Sprintf("webhook provider returned HTTP %d", err.statusCode)
}

func (err webhookResponseError) HTTPStatusCode() int {
	return err.statusCode
}

func (webhook *Webhook) SendWithContext(ctx context.Context, message message.Message) error {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.options.Url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return webhookResponseError{statusCode: resp.StatusCode}
	}
	return nil
}

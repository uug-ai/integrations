package integrations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// WebhookHTTPClient is an interface for sending HTTP requests
type WebhookHTTPClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// WebhookHTTPClientImpl is the default implementation using http.Client
type WebhookHTTPClientImpl struct {
	client *http.Client
}

// Post implements WebhookHTTPClient interface
func (w *WebhookHTTPClientImpl) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return w.client.Post(url, contentType, body)
}

// NewWebhookHTTPClient creates a new default webhook HTTP client with the specified timeout
func NewWebhookHTTPClient(timeout time.Duration) WebhookHTTPClient {
	return &WebhookHTTPClientImpl{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// WebhookOptions holds the configuration for Webhook
type WebhookOptions struct {
	Url     string `validate:"required,url"`
	Timeout int    `validate:"omitempty,gt=0"`
}

// WebhookOptionsBuilder provides a fluent interface for building Webhook options
type WebhookOptionsBuilder struct {
	options *WebhookOptions
}

// NewWebhookOptions creates a new Webhook options builder
func NewWebhookOptions() *WebhookOptionsBuilder {
	return &WebhookOptionsBuilder{
		options: &WebhookOptions{},
	}
}

// SetUrl sets the webhook URL
func (b *WebhookOptionsBuilder) SetUrl(url string) *WebhookOptionsBuilder {
	b.options.Url = url
	return b
}

// SetTimeout sets the timeout in seconds
func (b *WebhookOptionsBuilder) SetTimeout(timeout int) *WebhookOptionsBuilder {
	b.options.Timeout = timeout
	return b
}

// Build returns the configured WebhookOptions
func (b *WebhookOptionsBuilder) Build() *WebhookOptions {
	return b.options
}

// Webhook represents a Webhook client instance
type Webhook struct {
	options *WebhookOptions
	client  WebhookHTTPClient
}

// NewWebhook creates a new Webhook client with the provided options
// If client is not provided, a default WebhookHTTPClient will be created
func NewWebhook(opts *WebhookOptions, client ...WebhookHTTPClient) (*Webhook, error) {
	// Validate Webhook configuration
	validate := validator.New()
	err := validate.Struct(opts)
	if err != nil {
		return nil, err
	}

	// If no client provided, create default production client
	var c WebhookHTTPClient
	if len(client) == 0 || client[0] == nil {
		timeout := time.Duration(5 * time.Second)
		if opts.Timeout > 0 {
			timeout = time.Duration(opts.Timeout) * time.Second
		}
		c = NewWebhookHTTPClient(timeout)
	} else {
		c = client[0]
	}

	return &Webhook{
		options: opts,
		client:  c,
	}, nil
}

// Send sends a JSON payload to the webhook URL
// Parameters:
//   - body: The message or data to send as JSON
//
// Returns:
//   - error: An error if body is empty, if JSON marshaling fails, or if the HTTP request fails
func (w *Webhook) Send(body string) error {
	if body == "" {
		return errors.New("message body is empty")
	}
	// Prepare payload
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// Send HTTP POST request to the webhook URL
	resp, err := w.client.Post(w.options.Url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check if the request was successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook request failed with status: %s", resp.Status)
	}
	return nil
}

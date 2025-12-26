package integrations

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

// MockWebhookHTTPClient is a mock implementation of WebhookHTTPClient for testing
type MockWebhookHTTPClient struct {
	PostFunc     func(url string, contentType string, body io.Reader) (*http.Response, error)
	PostCalled   bool
	LastURL      string
	LastBodyType string
	LastBody     string
}

func (m *MockWebhookHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	m.PostCalled = true
	m.LastURL = url
	m.LastBodyType = contentType
	if body != nil {
		bodyBytes, _ := io.ReadAll(body)
		m.LastBody = string(bodyBytes)
	}
	if m.PostFunc != nil {
		return m.PostFunc(url, contentType, body)
	}
	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func setupWebhookTest() (*Webhook, error) {
	timeout, _ := strconv.Atoi(os.Getenv("WEBHOOK_TIMEOUT"))
	opts := NewWebhookOptions().
		SetUrl(os.Getenv("WEBHOOK_URL")).
		SetTimeout(timeout).
		Build()

	webhook, err := NewWebhook(opts, nil) // nil uses default production client
	return webhook, err
}

func TestWebhookValidation(t *testing.T) {
	// Use mock client to avoid actual network calls
	mockClient := &MockWebhookHTTPClient{}

	tests := []struct {
		name        string
		buildOpts   func() *WebhookOptions
		expectError bool
	}{
		{
			name: "MissingUrl",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetTimeout(5).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidUrl",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetUrl("not-a-valid-url").
					SetTimeout(5).
					Build()
			},
			expectError: true,
		},
		{
			name: "InvalidTimeout",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					SetTimeout(-1).
					Build()
			},
			expectError: true,
		},
		{
			name: "ZeroTimeout",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					SetTimeout(0).
					Build()
			},
			expectError: false, // 0 timeout is valid, will use default
		},
		{
			name: "ValidWithTimeout",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					SetTimeout(10).
					Build()
			},
			expectError: false,
		},
		{
			name: "ValidWithoutTimeout",
			buildOpts: func() *WebhookOptions {
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					Build()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()
			_, err := NewWebhook(opts, mockClient)
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

func TestWebhookFieldEmpty(t *testing.T) {
	// Use mock client to avoid actual network calls
	mockClient := &MockWebhookHTTPClient{}

	opts := NewWebhookOptions().
		SetUrl(os.Getenv("WEBHOOK_URL")).
		SetTimeout(5).
		Build()

	webhookIntegration, err := NewWebhook(opts, mockClient)
	if err != nil {
		t.Fatalf("failed to setup Webhook: %v", err)
	}

	tests := []struct {
		body        string
		expectError bool
	}{
		{body: "", expectError: true},
		{body: "Test message", expectError: false},
		{body: "{\"key\":\"value\"}", expectError: false},
	}

	for _, tt := range tests {
		mockClient.PostCalled = false // Reset for each test
		err := webhookIntegration.Send(tt.body)
		if tt.expectError && err == nil {
			t.Errorf("expected error got nil for body: '%s'", tt.body)
		}
		if !tt.expectError && err != nil {
			t.Errorf("expected no error got %v for body: '%s'", err, tt.body)
		}
		if !tt.expectError && !mockClient.PostCalled {
			t.Errorf("expected Post to be called but it wasn't")
		}
	}
}

func TestWebhookSendWithMock(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		mockPostError error
		mockStatus    int
		expectError   bool
		expectPost    bool
	}{
		{
			name:          "ValidSend",
			body:          "Test message",
			mockPostError: nil,
			mockStatus:    200,
			expectError:   false,
			expectPost:    true,
		},
		{
			name:          "EmptyBody",
			body:          "",
			mockPostError: nil,
			mockStatus:    200,
			expectError:   true,
			expectPost:    false,
		},
		{
			name:          "NetworkError",
			body:          "Test message",
			mockPostError: errors.New("network error"),
			mockStatus:    0,
			expectError:   true,
			expectPost:    true,
		},
		{
			name:          "BadRequest",
			body:          "Test message",
			mockPostError: nil,
			mockStatus:    400,
			expectError:   true,
			expectPost:    true,
		},
		{
			name:          "ServerError",
			body:          "Test message",
			mockPostError: nil,
			mockStatus:    500,
			expectError:   true,
			expectPost:    true,
		},
		{
			name:          "SuccessCreated",
			body:          "Test message",
			mockPostError: nil,
			mockStatus:    201,
			expectError:   false,
			expectPost:    true,
		},
		{
			name:          "SuccessNoContent",
			body:          "Test message",
			mockPostError: nil,
			mockStatus:    204,
			expectError:   false,
			expectPost:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client with specific error behavior
			mockClient := &MockWebhookHTTPClient{
				PostFunc: func(url string, contentType string, body io.Reader) (*http.Response, error) {
					if tt.mockPostError != nil {
						return nil, tt.mockPostError
					}
					return &http.Response{
						StatusCode: tt.mockStatus,
						Status:     http.StatusText(tt.mockStatus),
						Body:       io.NopCloser(strings.NewReader("")),
					}, nil
				},
			}

			opts := NewWebhookOptions().
				SetUrl(os.Getenv("WEBHOOK_URL")).
				SetTimeout(5).
				Build()

			webhookIntegration, err := NewWebhook(opts, mockClient)
			if err != nil {
				t.Fatalf("failed to setup Webhook: %v", err)
			}

			// Send message
			err = webhookIntegration.Send(tt.body)
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if tt.expectPost && !mockClient.PostCalled {
				t.Errorf("expected Post to be called but it wasn't")
			}
			if !tt.expectPost && mockClient.PostCalled {
				t.Errorf("expected Post not to be called but it was")
			}
		})
	}
}

func TestWebhookContentType(t *testing.T) {
	mockClient := &MockWebhookHTTPClient{}

	opts := NewWebhookOptions().
		SetUrl(os.Getenv("WEBHOOK_URL")).
		Build()

	webhookIntegration, err := NewWebhook(opts, mockClient)
	if err != nil {
		t.Fatalf("failed to setup Webhook: %v", err)
	}

	err = webhookIntegration.Send("Test message")
	if err != nil {
		t.Fatalf("failed to send webhook: %v", err)
	}

	if mockClient.LastBodyType != "application/json" {
		t.Errorf("expected content type 'application/json', got '%s'", mockClient.LastBodyType)
	}
}

func TestWebhookURL(t *testing.T) {
	expectedURL := "https://example.com/webhook"
	mockClient := &MockWebhookHTTPClient{}

	opts := NewWebhookOptions().
		SetUrl(expectedURL).
		Build()

	webhookIntegration, err := NewWebhook(opts, mockClient)
	if err != nil {
		t.Fatalf("failed to setup Webhook: %v", err)
	}

	err = webhookIntegration.Send("Test message")
	if err != nil {
		t.Fatalf("failed to send webhook: %v", err)
	}

	if mockClient.LastURL != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, mockClient.LastURL)
	}
}

func TestWebhookBuilderPattern(t *testing.T) {
	mockClient := &MockWebhookHTTPClient{}

	// Test fluent builder interface
	opts := NewWebhookOptions().
		SetUrl("https://example.com/webhook").
		SetTimeout(10).
		Build()

	if opts.Url != "https://example.com/webhook" {
		t.Errorf("expected URL 'https://example.com/webhook', got '%s'", opts.Url)
	}
	if opts.Timeout != 10 {
		t.Errorf("expected timeout 10, got %d", opts.Timeout)
	}

	webhook, err := NewWebhook(opts, mockClient)
	if err != nil {
		t.Fatalf("failed to create webhook: %v", err)
	}

	if webhook.options.Url != opts.Url {
		t.Errorf("webhook URL mismatch")
	}
}

func TestIntegrationWebhook(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	tests := []struct {
		name        string
		buildOpts   func() *WebhookOptions
		body        string
		expectError bool
	}{
		{
			name: "ValidWebhook",
			buildOpts: func() *WebhookOptions {
				timeout, _ := strconv.Atoi(os.Getenv("WEBHOOK_TIMEOUT"))
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					SetTimeout(timeout).
					Build()
			},
			body:        "{\"message\": \"Test message from UUG AI\"}",
			expectError: false,
		},
		{
			name: "ValidWebhookWithJSON",
			buildOpts: func() *WebhookOptions {
				timeout, _ := strconv.Atoi(os.Getenv("WEBHOOK_TIMEOUT"))
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL")).
					SetTimeout(timeout).
					Build()
			},
			body:        "{\"message\": \"Test from UUG AI\", \"timestamp\": \"2023-01-01T00:00:00Z\"}",
			expectError: false,
		},
		{
			name: "ValidWebhookReturns400",
			buildOpts: func() *WebhookOptions {
				timeout, _ := strconv.Atoi(os.Getenv("WEBHOOK_TIMEOUT"))
				return NewWebhookOptions().
					SetUrl(os.Getenv("WEBHOOK_URL") + "?httpCode=400").
					SetTimeout(timeout).
					Build()
			},
			body:        "{\"message\": \"Test from UUG AI\", \"timestamp\": \"2023-01-01T00:00:00Z\"}",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.buildOpts()
			webhookIntegration, err := NewWebhook(opts, nil) // Use real HTTP client
			if err != nil {
				t.Fatalf("failed to setup Webhook: %v", err)
			}

			err = webhookIntegration.Send(tt.body)
			if tt.expectError && err == nil {
				t.Errorf("expected error got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
		})
	}
}

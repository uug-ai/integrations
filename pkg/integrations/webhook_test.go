package integrations

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/uug-ai/integrations/pkg/message"
)

func TestWebhookSendWithContextAcknowledgement(t *testing.T) {
	for _, status := range []int{200, 202, 204, 400, 403, 429, 500, 503} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.Method != http.MethodPost || request.Header.Get("Content-Type") != "application/json" {
					t.Error("webhook request lost method or content type")
				}
				writer.WriteHeader(status)
				if status != http.StatusNoContent {
					_, _ = io.WriteString(writer, "private response body")
				}
			}))
			defer server.Close()
			err := (&Webhook{options: &WebhookOptions{Url: server.URL + "/private-token"}}).SendWithContext(context.Background(), message.Message{})
			if status < 300 {
				if err != nil {
					t.Fatalf("successful acknowledgement failed: %v", err)
				}
				return
			}
			var response interface{ HTTPStatusCode() int }
			if !errors.As(err, &response) || response.HTTPStatusCode() != status {
				t.Fatalf("error = %v, want typed HTTP %d", err, status)
			}
			if strings.Contains(err.Error(), "private") || strings.Contains(err.Error(), server.URL) {
				t.Fatalf("response error exposed private data: %v", err)
			}
		})
	}
}

func TestWebhookSendWithContextPreservesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := (&Webhook{options: &WebhookOptions{Url: "http://127.0.0.1:1/private-token"}}).SendWithContext(ctx, message.Message{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want cancellation", err)
	}
}

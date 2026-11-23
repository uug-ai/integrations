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

type Webhook struct {
	Url string `json:"url,omitempty"`
}

func (webhook Webhook) Send(message message.Message) bool {
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

func (webhook Webhook) SendWithContext(ctx context.Context, message message.Message) error {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.Url, bytes.NewBuffer(bytesRepresentation))
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

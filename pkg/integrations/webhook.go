package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Webhook struct {
	Url     string `json:"url,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

func (webhook Webhook) Send(body string) error {
	// Prepare payload
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// Set timeout, default to 5 seconds if not specified
	timeout := time.Duration(5 * time.Second)
	if webhook.Timeout > 0 {
		timeout = time.Duration(time.Duration(webhook.Timeout) * time.Second)
	}
	// Send HTTP POST request to the webhook URL
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Post(webhook.Url, "application/json", bytes.NewBuffer(bytesRepresentation))
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

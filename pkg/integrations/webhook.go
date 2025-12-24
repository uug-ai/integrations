package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

type Webhook struct {
	Url     string `json:"url,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

func (webhook Webhook) Send(message models.Message) error {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}
	timeout := time.Duration(5 * time.Second)
	if webhook.Timeout > 0 {
		timeout = time.Duration(time.Duration(webhook.Timeout) * time.Second)
	}
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Post(webhook.Url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	} else {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		return nil
	}
}

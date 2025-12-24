package integrations

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

type Webhook struct {
	Url string `json:"url,omitempty"`
}

func (webhook Webhook) Send(message models.Message) bool {
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(webhook.Url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Print(err)
		return false
	} else {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		return true
	}
}

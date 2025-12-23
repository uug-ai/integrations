package channels

import (
	"bytes"
	"encoding/json"
	"net/http"

	message "github.com/uug-ai/hub-pipeline-notification/message"
)

type Alexa struct {
	AccessCode string `json:"accesscode,omitempty"`
}

func (alexa Alexa) Send(message message.Message) bool {

	url := "https://api.notifymyecho.com/v1/NotifyMe"
	notification := message.Body
	values := map[string]string{"notification": notification, "accessCode": alexa.AccessCode}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return true
}

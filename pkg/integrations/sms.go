package integrations

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/uug-ai/integrations/pkg/message"
)

const twilioAPIBase = "https://api.twilio.com/2010-04-01"

type Sms struct {
	AccountSID string `json:"accountsid,omitempty"`
	AuthToken  string `json:"authtoken,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

func (sms Sms) Send(m message.Message) bool {

	accountSid := sms.AccountSID
	authToken := sms.AuthToken

	from := sms.From
	to := sms.To
	message := "This is a test message"
	if err := sendTwilioSMS(twilioAPIBase, accountSid, authToken, from, to, message); err != nil {
		log.Print(err)
	}

	return true
}

// sendTwilioSMS posts a message to Twilio's Messages API. It replaces
// github.com/sfreiberg/gotwilio, which pulled in golang-jwt/jwt v3
// (GO-2025-3553) and gorilla/schema v1.1.0 (GO-2024-2958).
func sendTwilioSMS(apiBase, accountSid, authToken, from, to, body string) error {
	form := url.Values{}
	form.Set("From", from)
	form.Set("To", to)
	form.Set("Body", body)

	endpoint := apiBase + "/Accounts/" + url.PathEscape(accountSid) + "/Messages.json"
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("twilio: send sms failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

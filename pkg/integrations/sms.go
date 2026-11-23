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

// SmsOptions configures Sms.
type SmsOptions struct {
	// AccountSID is the Twilio account SID.
	AccountSID string `json:"accountsid,omitempty" validate:"required"`
	// AuthToken is the Twilio auth token.
	AuthToken string `json:"authtoken,omitempty" validate:"required"`
	// From is the sending phone number.
	From string `json:"from,omitempty" validate:"required"`
	// To is the receiving phone number.
	To string `json:"to,omitempty" validate:"required"`
}

// SmsOptionsBuilder builds SmsOptions.
type SmsOptionsBuilder struct {
	options *SmsOptions
}

// NewSmsOptions starts building SmsOptions.
func NewSmsOptions() *SmsOptionsBuilder {
	return &SmsOptionsBuilder{options: &SmsOptions{}}
}

// SetAccountSID sets the Twilio account SID.
func (b *SmsOptionsBuilder) SetAccountSID(accountSID string) *SmsOptionsBuilder {
	b.options.AccountSID = accountSID
	return b
}

// SetAuthToken sets the Twilio auth token.
func (b *SmsOptionsBuilder) SetAuthToken(authToken string) *SmsOptionsBuilder {
	b.options.AuthToken = authToken
	return b
}

// SetFrom sets the sending phone number.
func (b *SmsOptionsBuilder) SetFrom(from string) *SmsOptionsBuilder {
	b.options.From = from
	return b
}

// SetTo sets the receiving phone number.
func (b *SmsOptionsBuilder) SetTo(to string) *SmsOptionsBuilder {
	b.options.To = to
	return b
}

// Build returns the configured SmsOptions.
func (b *SmsOptionsBuilder) Build() *SmsOptions {
	return b.options
}

// Sms sends a text message through Twilio.
type Sms struct {
	options *SmsOptions
}

// NewSms validates opts and returns a Sms sender that uses a copy of them.
func NewSms(opts *SmsOptions) (*Sms, error) {
	if err := validateOptions("sms", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Sms{options: &options}, nil
}

func (sms *Sms) Send(m message.Message) bool {

	accountSid := sms.options.AccountSID
	authToken := sms.options.AuthToken

	from := sms.options.From
	to := sms.options.To
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

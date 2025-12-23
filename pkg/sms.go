package pkg

import (
	"github.com/sfreiberg/gotwilio"
)

type Sms struct {
	AccountSID string `json:"accountsid,omitempty"`
	AuthToken  string `json:"authtoken,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

func (sms Sms) Send(m Message) bool {

	accountSid := sms.AccountSID
	authToken := sms.AuthToken
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := sms.From
	to := sms.To
	message := "This is a test message"
	twilio.SendSMS(from, to, message, "", "")

	return true
}

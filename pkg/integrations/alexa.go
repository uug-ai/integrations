package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/uug-ai/integrations/pkg/message"
)

// AlexaOptions configures Alexa.
type AlexaOptions struct {
	// AccessCode is the Notify Me access code.
	AccessCode string `json:"accesscode,omitempty" validate:"required"`
}

// AlexaOptionsBuilder builds AlexaOptions.
type AlexaOptionsBuilder struct {
	options *AlexaOptions
}

// NewAlexaOptions starts building AlexaOptions.
func NewAlexaOptions() *AlexaOptionsBuilder {
	return &AlexaOptionsBuilder{options: &AlexaOptions{}}
}

// SetAccessCode sets the Notify Me access code.
func (b *AlexaOptionsBuilder) SetAccessCode(accessCode string) *AlexaOptionsBuilder {
	b.options.AccessCode = accessCode
	return b
}

// Build returns the configured AlexaOptions.
func (b *AlexaOptionsBuilder) Build() *AlexaOptions {
	return b.options
}

// Alexa sends the message to an Echo through Notify Me.
type Alexa struct {
	options *AlexaOptions
}

// NewAlexa validates opts and returns a Alexa sender that uses a copy of them.
func NewAlexa(opts *AlexaOptions) (*Alexa, error) {
	if err := validateOptions("alexa", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Alexa{options: &options}, nil
}

func (alexa *Alexa) Send(message message.Message) bool {

	url := "https://api.notifymyecho.com/v1/NotifyMe"
	notification := message.Body
	values := map[string]string{"notification": notification, "accessCode": alexa.options.AccessCode}
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

package integrations

import (
	"fmt"

	pusho "github.com/gregdel/pushover"
	"github.com/uug-ai/integrations/pkg/message"
)

// PushoverOptions configures Pushover.
type PushoverOptions struct {
	// ApiKey is the Pushover application token.
	ApiKey string `json:"api_key,omitempty" validate:"required"`
	// SendTo is the user or group key.
	SendTo string `json:"send_to,omitempty" validate:"required"`
}

// PushoverOptionsBuilder builds PushoverOptions.
type PushoverOptionsBuilder struct {
	options *PushoverOptions
}

// NewPushoverOptions starts building PushoverOptions.
func NewPushoverOptions() *PushoverOptionsBuilder {
	return &PushoverOptionsBuilder{options: &PushoverOptions{}}
}

// SetApiKey sets the Pushover application token.
func (b *PushoverOptionsBuilder) SetApiKey(apiKey string) *PushoverOptionsBuilder {
	b.options.ApiKey = apiKey
	return b
}

// SetSendTo sets the user or group key.
func (b *PushoverOptionsBuilder) SetSendTo(sendTo string) *PushoverOptionsBuilder {
	b.options.SendTo = sendTo
	return b
}

// Build returns the configured PushoverOptions.
func (b *PushoverOptionsBuilder) Build() *PushoverOptions {
	return b.options
}

// Pushover sends the message as a Pushover notification.
type Pushover struct {
	options *PushoverOptions
}

// NewPushover validates opts and returns a Pushover sender that uses a copy of them.
func NewPushover(opts *PushoverOptions) (*Pushover, error) {
	if err := validateOptions("pushover", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Pushover{options: &options}, nil
}

func (pushover *Pushover) Send(m message.Message) bool {

	// Create a new pushover app with a token
	app := pusho.New(pushover.options.ApiKey)

	// Create a new recipient
	recipient := pusho.NewRecipient(pushover.options.SendTo)

	// Create the message to send
	mess := pusho.NewMessage(m.Title + " " + m.Body)

	_, err := app.SendMessage(mess, recipient)
	if err != nil {
		fmt.Println(err)
	}

	return true
}

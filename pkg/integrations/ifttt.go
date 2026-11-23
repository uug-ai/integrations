package integrations

import (
	"strconv"

	ift "github.com/lorenzobenvenuti/ifttt"
	"github.com/uug-ai/integrations/pkg/message"
)

// IftttOptions configures Ifttt.
type IftttOptions struct {
	// Token is the IFTTT webhook key.
	Token string `json:"token,omitempty" validate:"required"`
}

// IftttOptionsBuilder builds IftttOptions.
type IftttOptionsBuilder struct {
	options *IftttOptions
}

// NewIftttOptions starts building IftttOptions.
func NewIftttOptions() *IftttOptionsBuilder {
	return &IftttOptionsBuilder{options: &IftttOptions{}}
}

// SetToken sets the IFTTT webhook key.
func (b *IftttOptionsBuilder) SetToken(token string) *IftttOptionsBuilder {
	b.options.Token = token
	return b
}

// Build returns the configured IftttOptions.
func (b *IftttOptionsBuilder) Build() *IftttOptions {
	return b.options
}

// Ifttt triggers an IFTTT webhook event named after the message type.
type Ifttt struct {
	options *IftttOptions
}

// NewIfttt validates opts and returns a Ifttt sender that uses a copy of them.
func NewIfttt(opts *IftttOptions) (*Ifttt, error) {
	if err := validateOptions("ifttt", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Ifttt{options: &options}, nil
}

func (ifthisthenthat *Ifttt) Send(m message.Message) bool {

	iftttClient := ift.NewIftttClient(ifthisthenthat.options.Token)
	values := []string{m.Title, m.Body, strconv.FormatInt(m.Timestamp, 10)}
	iftttClient.Trigger(m.Type, values)

	return true
}

package integrations

import (
	"github.com/uug-ai/integrations/pkg/message"
	pushb "github.com/xconstruct/go-pushbullet"
)

// PushbulletOptions configures Pushbullet.
type PushbulletOptions struct {
	// ApiKey is the Pushbullet access token.
	ApiKey string `json:"api_key,omitempty" validate:"required"`
}

// PushbulletOptionsBuilder builds PushbulletOptions.
type PushbulletOptionsBuilder struct {
	options *PushbulletOptions
}

// NewPushbulletOptions starts building PushbulletOptions.
func NewPushbulletOptions() *PushbulletOptionsBuilder {
	return &PushbulletOptionsBuilder{options: &PushbulletOptions{}}
}

// SetApiKey sets the Pushbullet access token.
func (b *PushbulletOptionsBuilder) SetApiKey(apiKey string) *PushbulletOptionsBuilder {
	b.options.ApiKey = apiKey
	return b
}

// Build returns the configured PushbulletOptions.
func (b *PushbulletOptionsBuilder) Build() *PushbulletOptions {
	return b.options
}

// Pushbullet pushes the message to every device of a Pushbullet account.
type Pushbullet struct {
	options *PushbulletOptions
}

// NewPushbullet validates opts and returns a Pushbullet sender that uses a copy of them.
func NewPushbullet(opts *PushbulletOptions) (*Pushbullet, error) {
	if err := validateOptions("pushbullet", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Pushbullet{options: &options}, nil
}

func (pushbullet *Pushbullet) SendLink(message message.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.options.ApiKey)

	// Determine the link to push; fall back to a note if no media is attached
	// so an empty Media slice never panics (which would lose the notification).
	link := ""
	if len(message.Media) > 0 {
		link = message.Media[0].Url
	}

	// Get all the devices
	devs, err := pb.Devices()
	if err == nil {
		// Send a message to the first device
		for _, dev := range devs {
			if link != "" {
				pb.PushLink(dev.Iden, message.Title, link, message.Body)
			} else {
				pb.PushNote(dev.Iden, message.Title, message.Body)
			}
		}
	}

	return true
}

func (pushbullet *Pushbullet) SendMessage(message message.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.options.ApiKey)

	// Get all the devices
	devs, err := pb.Devices()
	if err == nil {
		// Send a message to the first device
		for _, dev := range devs {
			pb.PushNote(dev.Iden, message.Title, message.Body)
		}
	}

	return true
}

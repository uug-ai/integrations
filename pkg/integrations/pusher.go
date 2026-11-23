package integrations

import (
	"fmt"
	"log"
	"os"
	"regexp"

	push "github.com/pusher/pusher-http-go"
	"github.com/uug-ai/integrations/pkg/message"
)

// pusherChannelPattern matches the characters Pusher allows in a channel name.
// See https://pusher.com/docs/channels/using_channels/channels/#channel-naming-conventions
var pusherChannelPattern = regexp.MustCompile(`^[a-zA-Z0-9_=@,.;-]+$`)

// validPusherChannel reports whether recipient is usable as a Pusher channel
// name. message.User is the ROUTING KEY here (it must equal the user's
// username, which the frontend subscribes on) and NOT a display label. A
// display name like "Cedric Verstraeten" contains a space and is rejected by
// Pusher, which silently drops the browser notification. Guarding here turns
// that silent failure into a logged error.
func validPusherChannel(recipient string) bool {
	return recipient != "" && pusherChannelPattern.MatchString(recipient)
}

// Following message is send to pusher, structure may be modified,
// if this is required by the web interface.

/*{
  "sequence": {
    "title": "Hey, something happened!",
    "text": "Activity was detected at your frontdoor on 14:53:49.",
    "images": [
      {
        "title": "14:53:49",
        "media": "https://kerberosaccept.s3.eu-west-1.amazonaws.com/gvelim/1514991229_6-781657_frontdoor_27-238-373-314_131_764.mp4?X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIE5D4XAKE2O7DT4A%2F20180103%2Feu-west-1%2Fs3%2Faws4_request&X-Amz-Date=20180103T145434Z&X-Amz-SignedHeaders=host&X-Amz-Expires=1800&X-Amz-Signature=a898db92bb1df40a6a7d08c04b719f53869e05979350e02277cc2a1d12ace7fa",
        "type": "video"
      },
      {
        "title": "14:53:49",
        "media": "https://kerberosaccept.s3.eu-west-1.amazonaws.com/gvelim/1514991229_6-781657_frontdoor_27-238-373-314_131_764.mp4?X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIE5D4XAKE2O7DT4A%2F20180103%2Feu-west-1%2Fs3%2Faws4_request&X-Amz-Date=20180103T145434Z&X-Amz-SignedHeaders=host&X-Amz-Expires=1800&X-Amz-Signature=a898db92bb1df40a6a7d08c04b719f53869e05979350e02277cc2a1d12ace7fa",
        "type": "video"
      }
    ]
  }
}*/

type PusherMessageWrapper struct {
	Sequence PusherMessage `json:"sequence,omitempty"`
}

type PusherMessage struct {
	Title string        `json:"title,omitempty"`
	Text  string        `json:"text,omitempty"`
	Media []PusherMedia `json:"images,omitempty"`
}

type PusherMedia struct {
	Title string `json:"title,omitempty"`
	Media string `json:"media,omitempty"`
	Type  string `json:"type,omitempty"`
}

// PusherOptions configures Pusher.
type PusherOptions struct {
	// AppID is the Pusher app ID.
	AppID string `json:"app_id,omitempty" validate:"required"`
	// Key is the Pusher app key.
	Key string `json:"key,omitempty" validate:"required"`
	// Secret is the Pusher app secret.
	Secret string `json:"secret,omitempty" validate:"required"`
	// Cluster is the Pusher cluster.
	Cluster string `json:"cluster,omitempty" validate:"required"`
	// Channel is the event name.
	Channel string `json:"channel,omitempty" validate:"required"`
}

// PusherOptionsBuilder builds PusherOptions.
type PusherOptionsBuilder struct {
	options *PusherOptions
}

// NewPusherOptions starts building PusherOptions.
func NewPusherOptions() *PusherOptionsBuilder {
	return &PusherOptionsBuilder{options: &PusherOptions{}}
}

// SetAppID sets the Pusher app ID.
func (b *PusherOptionsBuilder) SetAppID(appID string) *PusherOptionsBuilder {
	b.options.AppID = appID
	return b
}

// SetKey sets the Pusher app key.
func (b *PusherOptionsBuilder) SetKey(key string) *PusherOptionsBuilder {
	b.options.Key = key
	return b
}

// SetSecret sets the Pusher app secret.
func (b *PusherOptionsBuilder) SetSecret(secret string) *PusherOptionsBuilder {
	b.options.Secret = secret
	return b
}

// SetCluster sets the Pusher cluster.
func (b *PusherOptionsBuilder) SetCluster(cluster string) *PusherOptionsBuilder {
	b.options.Cluster = cluster
	return b
}

// SetChannel sets the event name.
func (b *PusherOptionsBuilder) SetChannel(channel string) *PusherOptionsBuilder {
	b.options.Channel = channel
	return b
}

// Build returns the configured PusherOptions.
func (b *PusherOptionsBuilder) Build() *PusherOptions {
	return b.options
}

// Pusher triggers an event on the recipient's Pusher channel (Message.User).
type Pusher struct {
	options *PusherOptions
}

// NewPusher validates opts and returns a Pusher sender that uses a copy of them.
func NewPusher(opts *PusherOptions) (*Pusher, error) {
	if err := validateOptions("pusher", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Pusher{options: &options}, nil
}

// PusherFromEnv builds a Pusher sender for the given event name, using
// PUSHER_APP_ID, PUSHER_KEY, PUSHER_SECRET and PUSHER_CLUSTER.
func PusherFromEnv(channel string) (*Pusher, error) {
	return NewPusher(NewPusherOptions().
		SetAppID(os.Getenv("PUSHER_APP_ID")).
		SetKey(os.Getenv("PUSHER_KEY")).
		SetSecret(os.Getenv("PUSHER_SECRET")).
		SetCluster(os.Getenv("PUSHER_CLUSTER")).
		SetChannel(channel).
		Build())
}

func (pusher *Pusher) client() push.Client {
	return push.Client{
		AppID:   pusher.options.AppID,
		Key:     pusher.options.Key,
		Secret:  pusher.options.Secret,
		Cluster: pusher.options.Cluster,
	}
}

// buildPusherPayload assembles the wrapper sent to Pusher. It is a pure
// function (no network) so the payload mapping can be regression-tested, and it
// guards against an empty Media slice — a missing guard there would panic and,
// via the handler's recover, dead-letter the event (i.e. lose the in-app
// notification entirely).
func buildPusherPayload(message message.Message) PusherMessageWrapper {
	pusherMessage := PusherMessageWrapper{}
	pusherMessage.Sequence.Title = message.Title
	pusherMessage.Sequence.Text = message.Body
	pusherMessage.Sequence.Media = []PusherMedia{}
	if len(message.Media) > 0 {
		pusherMessage.Sequence.Media = append(pusherMessage.Sequence.Media, PusherMedia{
			Title: fmt.Sprintf("%v", message.Media[0].Timestamp),
			Type:  message.Media[0].Type,
			Media: message.Media[0].Url,
		})
	}
	return pusherMessage
}

func (pusher *Pusher) SendNotification(message message.Message) bool {

	if !validPusherChannel(message.User) {
		log.Printf("pusher: invalid channel %q (must be the user's username, not a display name); dropping notification", message.User)
		return false
	}

	client := pusher.client()

	pusherMessage := buildPusherPayload(message)

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.options.Channel, pusherMessage)

	return true
}

// Trigger sends the message on the recipient's channel (message.User) and
// reports why it was not delivered.
func (pusher *Pusher) Trigger(message message.Message) error {
	if !validPusherChannel(message.User) {
		return fmt.Errorf("pusher: invalid channel %q (must be the user's username, not a display name)", message.User)
	}
	client := pusher.client()
	return client.Trigger(message.User, pusher.options.Channel, message)
}

func (pusher *Pusher) Send(message message.Message) bool {

	if !validPusherChannel(message.User) {
		log.Printf("pusher: invalid channel %q (must be the user's username, not a display name); dropping notification", message.User)
		return false
	}

	client := pusher.client()

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.options.Channel, message)

	return true
}

package integrations

import (
	"fmt"
	"log"
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

// Pusher triggers an event on the recipient's Pusher channel. AppID, Key,
// Secret and Cluster identify the Pusher app; Channel is the event name.
type Pusher struct {
	AppID   string `json:"app_id,omitempty"`
	Key     string `json:"key,omitempty"`
	Secret  string `json:"secret,omitempty"`
	Cluster string `json:"cluster,omitempty"`
	Channel string `json:"channel,omitempty"`
}

func (pusher Pusher) client() push.Client {
	return push.Client{
		AppID:   pusher.AppID,
		Key:     pusher.Key,
		Secret:  pusher.Secret,
		Cluster: pusher.Cluster,
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

func (pusher Pusher) SendNotification(message message.Message) bool {

	if !validPusherChannel(message.User) {
		log.Printf("pusher: invalid channel %q (must be the user's username, not a display name); dropping notification", message.User)
		return false
	}

	client := pusher.client()

	pusherMessage := buildPusherPayload(message)

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.Channel, pusherMessage)

	return true
}

// Trigger sends the message on the recipient's channel (message.User) and
// reports why it was not delivered.
func (pusher Pusher) Trigger(message message.Message) error {
	if !validPusherChannel(message.User) {
		return fmt.Errorf("pusher: invalid channel %q (must be the user's username, not a display name)", message.User)
	}
	client := pusher.client()
	return client.Trigger(message.User, pusher.Channel, message)
}

func (pusher Pusher) Send(message message.Message) bool {

	if !validPusherChannel(message.User) {
		log.Printf("pusher: invalid channel %q (must be the user's username, not a display name); dropping notification", message.User)
		return false
	}

	client := pusher.client()

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.Channel, message)

	return true
}

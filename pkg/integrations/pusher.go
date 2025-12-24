package integrations

import (
	"fmt"

	push "github.com/pusher/pusher-http-go"
	"github.com/uug-ai/models/pkg/models"
)

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

type Pusher struct {
	Channel string `json:"channel,omitempty"`
}

func (pusher Pusher) SendNotification(message models.Message) bool {

	// instantiate a client
	client := push.Client{
		AppID:   "256802",
		Key:     "dbfbe47444eddb7e21e5",
		Secret:  "57e0315f3e4246225e62",
		Cluster: "eu",
	}

	pusherMessage := PusherMessageWrapper{}
	pusherMessage.Sequence.Title = message.Title
	pusherMessage.Sequence.Text = message.Body
	pusherMessage.Sequence.Media = []PusherMedia{}
	pusherMessage.Sequence.Media = append(pusherMessage.Sequence.Media, PusherMedia{
		Title: fmt.Sprintf("%v", message.Media[0].StartTimestamp),
		Media: message.Media[0].AtRuntimeMetadata.VideoUrl,
	})

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.Channel, pusherMessage)

	return true
}

func (pusher Pusher) Send(message models.Message) bool {

	// instantiate a client
	client := push.Client{
		AppID:   "256802",
		Key:     "dbfbe47444eddb7e21e5",
		Secret:  "57e0315f3e4246225e62",
		Cluster: "eu",
	}

	// trigger an event on the users channel, along with a data payload.
	client.Trigger(message.User, pusher.Channel, message)

	return true
}

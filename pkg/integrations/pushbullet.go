package integrations

import (
	"github.com/uug-ai/integrations/pkg/message"
	pushb "github.com/xconstruct/go-pushbullet"
)

type Pushbullet struct {
	ApiKey string `json:"api_key,omitempty"`
}

func (pushbullet Pushbullet) SendLink(message message.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.ApiKey)

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

func (pushbullet Pushbullet) SendMessage(message message.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.ApiKey)

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

package integrations

import (
	"github.com/uug-ai/models/pkg/models"
	pushb "github.com/xconstruct/go-pushbullet"
)

type Pushbullet struct {
	ApiKey string `json:"api_key,omitempty"`
}

func (pushbullet Pushbullet) SendLink(message models.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.ApiKey)

	// Get all the devices
	devs, err := pb.Devices()
	if err == nil {
		// Send a message to the first device
		for _, dev := range devs {
			pb.PushLink(dev.Iden, message.Title, message.Media[0].AtRuntimeMetadata.VideoUrl, message.Body)
		}
	}

	return true
}

func (pushbullet Pushbullet) SendMessage(message models.Message) bool {

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

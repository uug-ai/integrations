package channels

import (
	message "github.com/uug-ai/hub-pipeline-notification/message"
	pushb "github.com/xconstruct/go-pushbullet"
)

type Pushbullet struct {
	ApiKey string `json:"api_key,omitempty"`
}

func (pushbullet Pushbullet) SendLink(message message.Message) bool {

	// Instantiate a client
	pb := pushb.New(pushbullet.ApiKey)

	// Get all the devices
	devs, err := pb.Devices()
	if err == nil {
		// Send a message to the first device
		for _, dev := range devs {
			pb.PushLink(dev.Iden, message.Title, message.Media[0].Url, message.Body)
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

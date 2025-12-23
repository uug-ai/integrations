package channels

import (
	"fmt"

	pusho "github.com/gregdel/pushover"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

type Pushover struct {
	ApiKey string `json:"api_key,omitempty"`
	SendTo string `json:"send_to,omitempty"`
}

func (pushover Pushover) Send(m message.Message) bool {

	// Create a new pushover app with a token
	app := pusho.New(pushover.ApiKey)

	// Create a new recipient
	recipient := pusho.NewRecipient(pushover.SendTo)

	// Create the message to send
	mess := pusho.NewMessage(m.Title + " " + m.Body)

	_, err := app.SendMessage(mess, recipient)
	if err != nil {
		fmt.Println(err)
	}

	return true
}

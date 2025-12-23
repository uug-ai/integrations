package integrations

import (
	"fmt"

	pusho "github.com/gregdel/pushover"
	"github.com/uug-ai/models/pkg/models"
)

type Pushover struct {
	ApiKey string `json:"api_key,omitempty"`
	SendTo string `json:"send_to,omitempty"`
}

func (pushover Pushover) Send(m models.Message) bool {
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

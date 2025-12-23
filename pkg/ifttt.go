package pkg

import (
	"strconv"

	ift "github.com/lorenzobenvenuti/ifttt"
)

type Ifttt struct {
	Token string `json:"token,omitempty"`
}

func (ifthisthenthat Ifttt) Send(m Message) bool {

	iftttClient := ift.NewIftttClient(ifthisthenthat.Token)
	values := []string{m.Title, m.Body, strconv.FormatInt(m.Timestamp, 10)}
	iftttClient.Trigger(m.Type, values)

	return true
}

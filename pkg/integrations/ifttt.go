package integrations

import (
	"strconv"

	ift "github.com/lorenzobenvenuti/ifttt"
	"github.com/uug-ai/models/pkg/models"
)

type Ifttt struct {
	Token string `json:"token,omitempty"`
}

func (ifthisthenthat Ifttt) Send(m models.Message) bool {

	iftttClient := ift.NewIftttClient(ifthisthenthat.Token)
	values := []string{m.Title, m.Body, strconv.FormatInt(m.Timestamp, 10)}
	iftttClient.Trigger(m.Type, values)

	return true
}

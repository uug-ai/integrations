package channels

import (
	"context"
	"time"

	//shorturl "github.com/subosito/shorturl"
	mailgun "github.com/mailgun/mailgun-go/v4"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

type Mail struct {
	Domain     string `json:"domain,omitempty"`
	ApiKey     string `json:"api_key,omitempty"`
	TemplateId string `json:"templateId,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	EmailFrom  string `json:"email_from,omitempty"`
}

func (mail Mail) Send(message message.Message) error {

	domain := mail.Domain
	ApiKey := mail.ApiKey

	mg := mailgun.NewMailgun(domain, ApiKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	msg := mg.NewMessage(mail.EmailFrom,
		message.Title,
		"")
	msg.SetTemplate(mail.TemplateId)

	// Add recipients
	if mail.EmailTo != "" {
		msg.AddRecipient(mail.EmailTo)
	} else {
		msg.AddRecipient(message.Email)
	}

	if message.NumberOfMedia != "" {
		msg.AddVariable("numberOfMedia", message.NumberOfMedia)
	}
	if message.DataUsage != "" {
		msg.AddVariable("dataUsage", message.DataUsage)
	}

	// Add the variables to be used by the template
	msg.AddVariable("user", message.User)
	msg.AddVariable("text", message.Body)

	if len(message.Media) > 0 {
		longUrl := message.Media[0].Url
		msg.AddVariable("link", longUrl)
	}

	// Iterate over data object and modify
	for key, element := range message.Data {
		msg.AddVariable(key, element)
	}

	_, _, err := mg.Send(ctx, msg)

	return err
}

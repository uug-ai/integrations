package pkg

import (
	sg "github.com/sendgrid/sendgrid-go"
	mail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sendgrid struct {
	ApiKey           string `json:"api_key,omitempty"`
	EmailFrom        string `json:"email_from,omitempty"`
	EmailFromDisplay string `json:"email_from_display,omitempty"`
	EmailTo          string `json:"email_to,omitempty"`
	TemplateId       string `json:"templateId,omitempty"`
}

func (sendg Sendgrid) SendNotification(message Message) bool {

	from := mail.NewEmail(sendg.EmailFromDisplay, sendg.EmailFrom)
	to := mail.NewEmail(message.User, message.Email)
	if sendg.EmailTo != "" {
		to = mail.NewEmail(sendg.EmailTo, sendg.EmailTo)
	}
	subject := message.Title
	content := mail.NewContent("text/html", "I'm replacing the <strong>body tag</strong>")
	m := mail.NewV3MailInit(from, subject, to, content)
	m.Personalizations[0].SetSubstitution("{{user}}", message.User)
	if message.NumberOfMedia != "" {
		m.Personalizations[0].SetSubstitution("{{numberOfMedia}}", message.NumberOfMedia)
	}
	if message.DataUsage != "" {
		m.Personalizations[0].SetSubstitution("{{dataUsage}}", message.DataUsage)
	}
	m.Personalizations[0].SetSubstitution("{{text}}", message.Body)

	if len(message.Media) > 0 {
		longUrl := message.Media[0].Url
		//provider := "tinyurl"
		url := longUrl
		/*provider := "tinyurl"
		shortenedUrl, err := shorturl.Shorten(longUrl, provider)
		if err == nil {
			url = string(shortenedUrl)
		}*/
		//if err == nil {
		m.Personalizations[0].SetSubstitution("{{link}}", string(url))
		//}
	}
	m.SetTemplateID(sendg.TemplateId)

	request := sg.GetRequest(sendg.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	sg.API(request)

	return true
}

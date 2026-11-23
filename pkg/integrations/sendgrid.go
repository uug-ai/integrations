package integrations

import (
	sg "github.com/sendgrid/sendgrid-go"
	mail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/uug-ai/integrations/pkg/message"
)

// SendgridOptions configures Sendgrid.
type SendgridOptions struct {
	// ApiKey is the SendGrid API key.
	ApiKey string `json:"api_key,omitempty" validate:"required"`
	// EmailFrom is the sender address.
	EmailFrom string `json:"email_from,omitempty" validate:"required"`
	// EmailFromDisplay is the sender name.
	EmailFromDisplay string `json:"email_from_display,omitempty"`
	// EmailTo is the recipient; empty sends to the message's email.
	EmailTo string `json:"email_to,omitempty"`
	// TemplateId is the SendGrid template.
	TemplateId string `json:"templateId,omitempty"`
}

// SendgridOptionsBuilder builds SendgridOptions.
type SendgridOptionsBuilder struct {
	options *SendgridOptions
}

// NewSendgridOptions starts building SendgridOptions.
func NewSendgridOptions() *SendgridOptionsBuilder {
	return &SendgridOptionsBuilder{options: &SendgridOptions{}}
}

// SetApiKey sets the SendGrid API key.
func (b *SendgridOptionsBuilder) SetApiKey(apiKey string) *SendgridOptionsBuilder {
	b.options.ApiKey = apiKey
	return b
}

// SetEmailFrom sets the sender address.
func (b *SendgridOptionsBuilder) SetEmailFrom(emailFrom string) *SendgridOptionsBuilder {
	b.options.EmailFrom = emailFrom
	return b
}

// SetEmailFromDisplay sets the sender name.
func (b *SendgridOptionsBuilder) SetEmailFromDisplay(emailFromDisplay string) *SendgridOptionsBuilder {
	b.options.EmailFromDisplay = emailFromDisplay
	return b
}

// SetEmailTo sets the recipient; empty sends to the message's email.
func (b *SendgridOptionsBuilder) SetEmailTo(emailTo string) *SendgridOptionsBuilder {
	b.options.EmailTo = emailTo
	return b
}

// SetTemplateId sets the SendGrid template.
func (b *SendgridOptionsBuilder) SetTemplateId(templateId string) *SendgridOptionsBuilder {
	b.options.TemplateId = templateId
	return b
}

// Build returns the configured SendgridOptions.
func (b *SendgridOptionsBuilder) Build() *SendgridOptions {
	return b.options
}

// Sendgrid sends the message through a SendGrid template.
type Sendgrid struct {
	options *SendgridOptions
}

// NewSendgrid validates opts and returns a Sendgrid sender that uses a copy of them.
func NewSendgrid(opts *SendgridOptions) (*Sendgrid, error) {
	if err := validateOptions("sendgrid", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Sendgrid{options: &options}, nil
}

func (sendg *Sendgrid) SendNotification(message message.Message) bool {

	from := mail.NewEmail(sendg.options.EmailFromDisplay, sendg.options.EmailFrom)
	to := mail.NewEmail(message.User, message.Email)
	if sendg.options.EmailTo != "" {
		to = mail.NewEmail(sendg.options.EmailTo, sendg.options.EmailTo)
	}
	subject := message.Title
	content := mail.NewContent("text/html", "I'm replacing the <strong>body tag</strong>")
	m := mail.NewV3MailInit(from, subject, to, content)
	displayUser := message.User
	if message.Data["user"] != "" {
		displayUser = message.Data["user"]
	}
	m.Personalizations[0].SetSubstitution("{{user}}", displayUser)
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
	m.SetTemplateID(sendg.options.TemplateId)

	request := sg.GetRequest(sendg.options.ApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	sg.API(request)

	return true
}

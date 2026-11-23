package integrations

import (
	"context"
	"time"

	//shorturl "github.com/subosito/shorturl"
	mailgun "github.com/mailgun/mailgun-go/v4"
	"github.com/uug-ai/integrations/pkg/message"
)

// MailOptions configures Mail.
type MailOptions struct {
	// Domain is the Mailgun domain.
	Domain string `json:"domain,omitempty" validate:"required"`
	// ApiKey is the Mailgun API key.
	ApiKey string `json:"api_key,omitempty" validate:"required"`
	// TemplateId is the Mailgun template.
	TemplateId string `json:"templateId,omitempty"`
	// EmailTo is the recipient; empty sends to the message's email.
	EmailTo string `json:"email_to,omitempty"`
	// EmailFrom is the sender address.
	EmailFrom string `json:"email_from,omitempty" validate:"required"`
}

// MailOptionsBuilder builds MailOptions.
type MailOptionsBuilder struct {
	options *MailOptions
}

// NewMailOptions starts building MailOptions.
func NewMailOptions() *MailOptionsBuilder {
	return &MailOptionsBuilder{options: &MailOptions{}}
}

// SetDomain sets the Mailgun domain.
func (b *MailOptionsBuilder) SetDomain(domain string) *MailOptionsBuilder {
	b.options.Domain = domain
	return b
}

// SetApiKey sets the Mailgun API key.
func (b *MailOptionsBuilder) SetApiKey(apiKey string) *MailOptionsBuilder {
	b.options.ApiKey = apiKey
	return b
}

// SetTemplateId sets the Mailgun template.
func (b *MailOptionsBuilder) SetTemplateId(templateId string) *MailOptionsBuilder {
	b.options.TemplateId = templateId
	return b
}

// SetEmailTo sets the recipient; empty sends to the message's email.
func (b *MailOptionsBuilder) SetEmailTo(emailTo string) *MailOptionsBuilder {
	b.options.EmailTo = emailTo
	return b
}

// SetEmailFrom sets the sender address.
func (b *MailOptionsBuilder) SetEmailFrom(emailFrom string) *MailOptionsBuilder {
	b.options.EmailFrom = emailFrom
	return b
}

// Build returns the configured MailOptions.
func (b *MailOptionsBuilder) Build() *MailOptions {
	return b.options
}

// Mail sends the message through a Mailgun template.
type Mail struct {
	options *MailOptions
}

// NewMail validates opts and returns a Mail sender that uses a copy of them.
func NewMail(opts *MailOptions) (*Mail, error) {
	if err := validateOptions("mailgun", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Mail{options: &options}, nil
}

func (mail *Mail) Send(message message.Message) error {

	domain := mail.options.Domain
	ApiKey := mail.options.ApiKey

	mg := mailgun.NewMailgun(domain, ApiKey)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	msg := mg.NewMessage(mail.options.EmailFrom,
		message.Title,
		"")
	msg.SetTemplate(mail.options.TemplateId)

	// Add recipients
	if mail.options.EmailTo != "" {
		msg.AddRecipient(mail.options.EmailTo)
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

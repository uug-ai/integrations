package integrations

import (
	"fmt"
	"os"

	"github.com/uug-ai/integrations/pkg/message"
)

// Mail providers, as named by the MAIL_PROVIDER setting.
const (
	EmailProviderSMTP     = "smtp"
	EmailProviderMailgun  = "mailgun"
	EmailProviderSendgrid = "sendgrid"
)

// Email sends a message through one mail provider.
type Email interface {
	Send(m message.Message) error
}

// UnsupportedEmailProviderError is returned for a mail provider this package
// does not know.
type UnsupportedEmailProviderError string

func (err UnsupportedEmailProviderError) Error() string {
	return fmt.Sprintf("email provider %q is not supported", string(err))
}

// SMTPFromEnv builds an SMTP sender to emailTo with the given template, using
// SMTP_SERVER, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD and EMAIL_FROM.
func SMTPFromEnv(emailTo, templateId string) SMTP {
	return SMTP{
		Server:     os.Getenv("SMTP_SERVER"),
		Port:       os.Getenv("SMTP_PORT"),
		Username:   os.Getenv("SMTP_USERNAME"),
		Password:   os.Getenv("SMTP_PASSWORD"),
		EmailFrom:  os.Getenv("EMAIL_FROM"),
		EmailTo:    emailTo,
		TemplateId: templateId,
	}
}

// MailgunFromEnv builds a Mailgun sender to emailTo with the given template,
// using MAILGUN_DOMAIN, MAILGUN_API_KEY and EMAIL_FROM.
func MailgunFromEnv(emailTo, templateId string) Mail {
	return Mail{
		Domain:     os.Getenv("MAILGUN_DOMAIN"),
		ApiKey:     os.Getenv("MAILGUN_API_KEY"),
		EmailFrom:  os.Getenv("EMAIL_FROM"),
		EmailTo:    emailTo,
		TemplateId: templateId,
	}
}

// SendgridFromEnv builds a SendGrid sender to emailTo with the given template,
// using SENDGRID_API_KEY, EMAIL_FROM and EMAIL_FROM_DISPLAYNAME.
func SendgridFromEnv(emailTo, templateId string) Sendgrid {
	return Sendgrid{
		ApiKey:           os.Getenv("SENDGRID_API_KEY"),
		EmailFrom:        os.Getenv("EMAIL_FROM"),
		EmailFromDisplay: os.Getenv("EMAIL_FROM_DISPLAYNAME"),
		EmailTo:          emailTo,
		TemplateId:       templateId,
	}
}

// EmailFromEnv builds the sender for the named mail provider (smtp, mailgun or
// sendgrid) from the environment, as the *FromEnv builders do.
func EmailFromEnv(provider, emailTo, templateId string) (Email, error) {
	switch provider {
	case EmailProviderSMTP:
		return SMTPFromEnv(emailTo, templateId), nil
	case EmailProviderMailgun:
		return MailgunFromEnv(emailTo, templateId), nil
	case EmailProviderSendgrid:
		return sendgridEmail{SendgridFromEnv(emailTo, templateId)}, nil
	}
	return nil, UnsupportedEmailProviderError(provider)
}

// sendgridEmail adapts Sendgrid, which does not report delivery errors, to Email.
type sendgridEmail struct {
	sendgrid Sendgrid
}

func (email sendgridEmail) Send(m message.Message) error {
	email.sendgrid.SendNotification(m)
	return nil
}

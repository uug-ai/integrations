package integrations

import (
	"fmt"
	"os"
	"strconv"

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

// SMTPFromEnv builds an SMTP sender to emailTo that renders the named template
// from templates, using SMTP_SERVER, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD
// and EMAIL_FROM.
func SMTPFromEnv(emailTo, templateId string, templates Templates) (*SMTP, error) {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return NewSMTP(NewSMTPOptions().
		SetServer(os.Getenv("SMTP_SERVER")).
		SetPort(port).
		SetUsername(os.Getenv("SMTP_USERNAME")).
		SetPassword(os.Getenv("SMTP_PASSWORD")).
		SetEmailFrom(os.Getenv("EMAIL_FROM")).
		SetEmailTo(emailTo).
		SetTemplateId(templateId).
		SetTemplates(templates).
		Build())
}

// MailgunFromEnv builds a Mailgun sender to emailTo with the given template,
// using MAILGUN_DOMAIN, MAILGUN_API_KEY and EMAIL_FROM.
func MailgunFromEnv(emailTo, templateId string) (*Mail, error) {
	return NewMail(NewMailOptions().
		SetDomain(os.Getenv("MAILGUN_DOMAIN")).
		SetApiKey(os.Getenv("MAILGUN_API_KEY")).
		SetEmailFrom(os.Getenv("EMAIL_FROM")).
		SetEmailTo(emailTo).
		SetTemplateId(templateId).
		Build())
}

// SendgridFromEnv builds a SendGrid sender to emailTo with the given template,
// using SENDGRID_API_KEY, EMAIL_FROM and EMAIL_FROM_DISPLAYNAME.
func SendgridFromEnv(emailTo, templateId string) (*Sendgrid, error) {
	return NewSendgrid(NewSendgridOptions().
		SetApiKey(os.Getenv("SENDGRID_API_KEY")).
		SetEmailFrom(os.Getenv("EMAIL_FROM")).
		SetEmailFromDisplay(os.Getenv("EMAIL_FROM_DISPLAYNAME")).
		SetEmailTo(emailTo).
		SetTemplateId(templateId).
		Build())
}

// EmailFromEnv builds the sender for the named mail provider (smtp, mailgun or
// sendgrid) from the environment, as the *FromEnv builders do. Only SMTP
// renders templates; Mailgun and SendGrid use the template stored under
// templateId at the provider.
func EmailFromEnv(provider, emailTo, templateId string, templates Templates) (Email, error) {
	switch provider {
	case EmailProviderSMTP:
		smtp, err := SMTPFromEnv(emailTo, templateId, templates)
		if err != nil {
			return nil, err
		}
		return smtp, nil
	case EmailProviderMailgun:
		mail, err := MailgunFromEnv(emailTo, templateId)
		if err != nil {
			return nil, err
		}
		return mail, nil
	case EmailProviderSendgrid:
		sendgrid, err := SendgridFromEnv(emailTo, templateId)
		if err != nil {
			return nil, err
		}
		return sendgridEmail{sendgrid}, nil
	}
	return nil, UnsupportedEmailProviderError(provider)
}

// sendgridEmail adapts Sendgrid, which does not report delivery errors, to Email.
type sendgridEmail struct {
	sendgrid *Sendgrid
}

func (email sendgridEmail) Send(m message.Message) error {
	email.sendgrid.SendNotification(m)
	return nil
}

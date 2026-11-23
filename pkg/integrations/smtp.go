package integrations

import (
	"bytes"
	"crypto/tls"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/uug-ai/integrations/pkg/message"
	"gopkg.in/gomail.v2"
)

// SMTPOptions configures SMTP.
type SMTPOptions struct {
	// Server is the SMTP server host.
	Server string `json:"server,omitempty" validate:"required"`
	// Port is the SMTP server port.
	Port int `json:"port,omitempty" validate:"required,gt=0,lte=65535"`
	// Username is the SMTP user.
	Username string `json:"username,omitempty"`
	// Password is the SMTP password.
	Password string `json:"password,omitempty"`
	// EmailFrom is the sender address.
	EmailFrom string `json:"email_from,omitempty" validate:"required"`
	// EmailTo is the recipient address.
	EmailTo string `json:"email_to,omitempty" validate:"required"`
	// TemplateId is the name of the template to render.
	TemplateId string `json:"template_id,omitempty"`
	// Templates is the email templates TemplateId names.
	Templates Templates `json:"-" validate:"-"`
}

// SMTPOptionsBuilder builds SMTPOptions.
type SMTPOptionsBuilder struct {
	options *SMTPOptions
}

// NewSMTPOptions starts building SMTPOptions.
func NewSMTPOptions() *SMTPOptionsBuilder {
	return &SMTPOptionsBuilder{options: &SMTPOptions{}}
}

// SetServer sets the SMTP server host.
func (b *SMTPOptionsBuilder) SetServer(server string) *SMTPOptionsBuilder {
	b.options.Server = server
	return b
}

// SetPort sets the SMTP server port.
func (b *SMTPOptionsBuilder) SetPort(port int) *SMTPOptionsBuilder {
	b.options.Port = port
	return b
}

// SetUsername sets the SMTP user.
func (b *SMTPOptionsBuilder) SetUsername(username string) *SMTPOptionsBuilder {
	b.options.Username = username
	return b
}

// SetPassword sets the SMTP password.
func (b *SMTPOptionsBuilder) SetPassword(password string) *SMTPOptionsBuilder {
	b.options.Password = password
	return b
}

// SetEmailFrom sets the sender address.
func (b *SMTPOptionsBuilder) SetEmailFrom(emailFrom string) *SMTPOptionsBuilder {
	b.options.EmailFrom = emailFrom
	return b
}

// SetEmailTo sets the recipient address.
func (b *SMTPOptionsBuilder) SetEmailTo(emailTo string) *SMTPOptionsBuilder {
	b.options.EmailTo = emailTo
	return b
}

// SetTemplateId sets the name of the template to render.
func (b *SMTPOptionsBuilder) SetTemplateId(templateId string) *SMTPOptionsBuilder {
	b.options.TemplateId = templateId
	return b
}

// SetTemplates sets the email templates TemplateId names.
func (b *SMTPOptionsBuilder) SetTemplates(templates Templates) *SMTPOptionsBuilder {
	b.options.Templates = templates
	return b
}

// Build returns the configured SMTPOptions.
func (b *SMTPOptionsBuilder) Build() *SMTPOptions {
	return b.options
}

// SMTP sends the message as an email rendered from a template.
type SMTP struct {
	options *SMTPOptions
}

// NewSMTP validates opts and returns a SMTP sender that uses a copy of them.
func NewSMTP(opts *SMTPOptions) (*SMTP, error) {
	if err := validateOptions("smtp", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &SMTP{options: &options}, nil
}

func (smtp *SMTP) Send(message message.Message) error {
	d := gomail.NewDialer(smtp.options.Server, smtp.options.Port, smtp.options.Username, smtp.options.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", smtp.options.EmailFrom)
	m.SetHeader("To", smtp.options.EmailTo)
	// Strip HTML tags (e.g. <br>) from title for the plain-text subject line.
	subject := strings.ReplaceAll(message.Title, "<br>", " ")
	m.SetHeader("Subject", subject)
	timeNow := time.Now().Unix()
	m.SetHeader("Message-Id", "<"+strconv.FormatInt(timeNow, 10)+"@uug.ai>")

	template := smtp.options.Templates.Get(smtp.options.TemplateId)
	m.SetBody("text/plain", ReplaceValues(template.Text, message))
	m.AddAlternative("text/html", ReplaceValues(template.HTML, message))

	// Embed thumbnail as CID attachment so Gmail and other clients display it reliably.
	if len(message.Thumbnail) > 0 {
		m.Embed("thumbnail.png", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, bytes.NewReader(message.Thumbnail))
			return err
		}))
	}

	return d.DialAndSend(m)
}

// This function will replace the variables in the email template. We have following variables available:
// - {{user}}: user that triggered the message
// - {{text}}: text of the message
// - {{link}}: link to the media (recording)
//
// - {{thumbnail}}: image (either a base64 or a url).
// - {{classifications}}: list of classifications detected in the recording.
//
// - {{timezone}}: timezone of the account generating the event
//
// - {{date}}: date of the media
// - {{time}}: time of the media
// - {{datetime}}: datetime of the media
//
// - {{eventdate}}: date of the notification
// - {{eventtime}}: time of the notification
// - {{eventdatetime}}: datetime of the notification
//
// - {{devicename}}: device generating the event
// - {{deviceid}}: device generating the event
// - {{sites}}: the list of sites the device is part of
// - {{groups}}: the list of groups the device is part of
// - {{numberOfMedia}}: number of media attached to the message
// - {{dataUsage}}: data usage of the message

func ReplaceValues(body string, message message.Message) string {

	body = strings.ReplaceAll(body, "{{tab1_title}}", "")
	body = strings.ReplaceAll(body, "{{tab2_title}}", "")

	// {{title}} this will inject the alert title into the header.
	if message.Title != "" {
		body = strings.ReplaceAll(body, "{{title}}", message.Title)
	}

	if message.NumberOfMedia != "" {
		body = strings.ReplaceAll(body, "{{numberOfMedia}}", message.NumberOfMedia)
	}
	if message.DataUsage != "" {
		body = strings.ReplaceAll(body, "{{dataUsage}}", message.DataUsage)
	}

	// Add the variables to be used by the template
	//body = strings.ReplaceAll(body, "{{user}}", "")
	if message.Data["user"] != "" {
		body = strings.ReplaceAll(body, "{{user}}", message.Data["user"])
	} else {
		body = strings.ReplaceAll(body, "{{user}}", message.User)
	}

	body = strings.ReplaceAll(body, "{{text}}", message.Body)

	// {{link}} this will inject a link to the media (recording)
	// Data["link"] (the tenant-aware media link) takes precedence over the raw
	// media URL, matching the Mailgun channel so the same alert links the same
	// place regardless of email provider.
	if message.Data["link"] != "" {
		body = strings.ReplaceAll(body, "{{link}}", message.Data["link"])
	} else if len(message.Media) > 0 {
		longUrl := message.Media[0].Url
		body = strings.ReplaceAll(body, "{{link}}", longUrl)
	}

	// {{medialink}} this will inject a link to the media page on app.kerberos.io
	if message.MediaId != "" {
		body = strings.ReplaceAll(body, "{{medialink}}", "https://app.kerberos.io/media/"+message.MediaId)
	}

	// {{thumbnail}} this will inject an image using CID embedding or a URL fallback.
	if len(message.Thumbnail) > 0 {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='cid:thumbnail.png' width='100%' height='auto' />")
	} else if len(message.Media) > 0 && message.Media[0].ThumbnailUrl != "" {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='"+message.Media[0].ThumbnailUrl+"' width='100%' height='auto' />")
	}

	// {{classifications}} this will inject a list of classifications detected in the recording.
	if len(message.Classifications) > 0 {
		classificationNames := []string{}
		for _, classification := range message.Classifications {
			classificationNames = append(classificationNames, classification.Key)
		}
		body = strings.ReplaceAll(body, "{{classifications}}", strings.Join(classificationNames, ", "))
	}

	// {{date}} {{time}} {{datetime}} of the start.
	if len(message.Media) > 0 && message.Media[0].Timestamp > 0 {
		t := time.Unix(message.Media[0].Timestamp, 0)
		// Get time with timezone
		if message.Timezone != "" {
			loc, _ := time.LoadLocation(message.Timezone)
			t = t.In(loc)
		}
		body = strings.ReplaceAll(body, "{{date}}", t.Format("2006-01-02"))
		body = strings.ReplaceAll(body, "{{time}}", t.Format("15:04:05"))
		body = strings.ReplaceAll(body, "{{datetime}}", t.Format("2006-01-02 15:04:05"))
	}

	// {{eventtime}} of the notification
	if message.Timestamp > 0 {
		t := time.Unix(message.Timestamp, 0)
		// Get time with timezone
		if message.Timezone != "" {
			loc, _ := time.LoadLocation(message.Timezone)
			t = t.In(loc)
		}
		body = strings.ReplaceAll(body, "{{eventdate}}", t.Format("2006-01-02"))
		body = strings.ReplaceAll(body, "{{eventtime}}", t.Format("15:04:05"))
		body = strings.ReplaceAll(body, "{{eventdatetime}}", t.Format("2006-01-02 15:04:05"))
	}

	// {{timezone}} of the account generating the event
	if message.Timezone != "" {
		body = strings.ReplaceAll(body, "{{timezone}}", message.Timezone)
	}

	// {{devicename}} {{deviceid}} device generating the event
	if message.DeviceId != "" {
		body = strings.ReplaceAll(body, "{{deviceid}}", message.DeviceId)
	}
	if message.DeviceName != "" {
		body = strings.ReplaceAll(body, "{{devicename}}", message.DeviceName)
	}

	// {{sites}} the list of sites the device is part of
	if len(message.Sites) > 0 {
		// get names of sites slice (extract name)
		sitesName := []string{}
		for _, site := range message.Sites {
			sitesName = append(sitesName, site.Name)
		}
		body = strings.ReplaceAll(body, "{{sites}}", strings.Join(sitesName, ", "))
	}

	// {{groups}} the list of groups the device is part of
	if len(message.Groups) > 0 {
		// get names of groups slice (extract name)
		groupsName := []string{}
		for _, group := range message.Groups {
			groupsName = append(groupsName, group.Name)
		}
		body = strings.ReplaceAll(body, "{{groups}}", strings.Join(groupsName, ", "))
	}

	// Whipe out all variables with the {{variable}} syntax (regex: {{.*}})
	body = strings.ReplaceAll(body, "{{.*}}", "")

	// Iterate over data object and modify
	for key, element := range message.Data {
		body = strings.ReplaceAll(body, "{{"+key+"}}", element)
	}

	return body
}

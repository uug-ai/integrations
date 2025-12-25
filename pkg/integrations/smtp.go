package integrations

import (
	"crypto/tls"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
)

// MailDialer is an interface for sending emails
type MailDialer interface {
	DialAndSend(m ...*gomail.Message) error
	Dial() (gomail.SendCloser, error)
}

// GomailDialer wraps gomail.Dialer to implement MailDialer interface
type GomailDialer struct {
	dialer *gomail.Dialer
}

// NewGomailDialer creates a new GomailDialer with the provided SMTP settings
func NewGomailDialer(host string, port int, username, password string) MailDialer {
	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &GomailDialer{dialer: d}
}

// DialAndSend implements MailDialer interface
func (g *GomailDialer) DialAndSend(m ...*gomail.Message) error {
	return g.dialer.DialAndSend(m...)
}

// Dial implements MailDialer interface
func (g *GomailDialer) Dial() (gomail.SendCloser, error) {
	return g.dialer.Dial()
}

// SMTPOptions holds the configuration for SMTP
type SMTPOptions struct {
	Server    string `validate:"required"`
	Port      int    `validate:"required,gt=0"`
	Username  string `validate:"required"`
	Password  string `validate:"required"`
	EmailFrom string `validate:"required,email"`
	EmailTo   string `validate:"required,email"`
}

// SMTPOptionsBuilder provides a fluent interface for building SMTP options
type SMTPOptionsBuilder struct {
	options *SMTPOptions
}

// SMTPOptions creates a new SMTP options builder
func NewSMTPOptions() *SMTPOptionsBuilder {
	return &SMTPOptionsBuilder{
		options: &SMTPOptions{},
	}
}

// SetServer sets the SMTP server hostname or IP address
func (b *SMTPOptionsBuilder) SetServer(server string) *SMTPOptionsBuilder {
	b.options.Server = server
	return b
}

// SetPort sets the SMTP server port
func (b *SMTPOptionsBuilder) SetPort(port int) *SMTPOptionsBuilder {
	b.options.Port = port
	return b
}

// SetUsername sets the SMTP authentication username
func (b *SMTPOptionsBuilder) SetUsername(username string) *SMTPOptionsBuilder {
	b.options.Username = username
	return b
}

// SetPassword sets the SMTP authentication password
func (b *SMTPOptionsBuilder) SetPassword(password string) *SMTPOptionsBuilder {
	b.options.Password = password
	return b
}

// SetFrom sets the sender email address
func (b *SMTPOptionsBuilder) SetFrom(emailFrom string) *SMTPOptionsBuilder {
	b.options.EmailFrom = emailFrom
	return b
}

// SetTo sets the recipient email address
func (b *SMTPOptionsBuilder) SetTo(emailTo string) *SMTPOptionsBuilder {
	b.options.EmailTo = emailTo
	return b
}

// Build returns the configured SMTPOptions
func (b *SMTPOptionsBuilder) Build() *SMTPOptions {
	return b.options
}

// SMTP represents an SMTP client instance
type SMTP struct {
	options *SMTPOptions
	dialer  MailDialer
}

// NewSMTP creates a new SMTP client with the provided options
// If dialer is not provided, a default GomailDialer will be created
func NewSMTP(opts *SMTPOptions, dialer ...MailDialer) (*SMTP, error) {
	// Validate SMTP configuration
	validate := validator.New()
	err := validate.Struct(opts)
	if err != nil {
		return nil, err
	}

	// If no dialer provided, create default production dialer
	var d MailDialer
	if len(dialer) == 0 {
		d = NewGomailDialer(opts.Server, opts.Port, opts.Username, opts.Password)
	} else {
		d = dialer[0]
	}

	return &SMTP{
		options: opts,
		dialer:  d,
	}, nil
}

func (s *SMTP) Send(title string, body string, textBody string) (err error) {

	// Check if title and body are not empty
	if title == "" {
		return errors.New("empty title")
	}
	if body == "" {
		return errors.New("empty body")
	}
	if textBody == "" {
		return errors.New("empty text body")
	}

	// Check if we can dial to the server
	_, err = s.dialer.Dial()
	if err != nil {
		return err
	}

	// Create the message
	m := gomail.NewMessage()
	m.SetHeader("From", s.options.EmailFrom)
	m.SetHeader("To", s.options.EmailTo)
	m.SetHeader("Subject", title)
	timeNow := time.Now().Unix()
	m.SetHeader("Message-Id", "<"+strconv.FormatInt(timeNow, 10)+"@kerberos.io>")

	// Replace needs to be moved outside and placed in the hub-pipeline-notification
	//body := templates.GetTemplate(smtp.TemplateId)
	//textBody := templates.GetTextTemplate(smtp.TemplateId)
	// Replace variables in the template
	//body = ReplaceValues(body, models.Message{})
	//textBody = ReplaceValues(textBody, models.Message{})

	m.SetBody("text/plain", body)
	m.AddAlternative("text/html", textBody)

	// Send the email
	err = s.dialer.DialAndSend(m)
	return err
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

/*func ReplaceValues(body string, message models.Message) string {

	body = strings.ReplaceAll(body, "{{tab1_title}}", "")
	body = strings.ReplaceAll(body, "{{tab2_title}}", "")

	if message.NumberOfMedia != "" {
		body = strings.ReplaceAll(body, "{{numberOfMedia}}", message.NumberOfMedia)
	}
	if message.DataUsage != "" {
		body = strings.ReplaceAll(body, "{{dataUsage}}", message.DataUsage)
	}

	// Add the variables to be used by the template
	//body = strings.ReplaceAll(body, "{{user}}", "")
	if message.User != "" {
		body = strings.ReplaceAll(body, "{{user}}", message.User)
	} else {
		body = strings.ReplaceAll(body, "{{user}}", message.Data["user"])
	}

	body = strings.ReplaceAll(body, "{{text}}", message.Body)

	// {{link}} this will inject a link to the media (recording)
	if len(message.Media) > 0 {
		longUrl := message.Media[0].AtRuntimeMetadata.VideoUrl
		body = strings.ReplaceAll(body, "{{link}}", longUrl)
	}
	if message.Data["link"] != "" {
		body = strings.ReplaceAll(body, "{{link}}", message.Data["link"])
	}

	// {{thumbnail}} this will inject an image (either a base64 or a url).
	if message.Thumbnail != "" {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='base64:"+message.Thumbnail+"' width='400px' height='auto' />")
	} else if len(message.Media) > 0 && message.Media[0].AtRuntimeMetadata.ThumbnailUrl != "" {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='"+message.Media[0].AtRuntimeMetadata.ThumbnailUrl+"' width='400px' height='auto' />")
	}

	// {{classifications}} this will inject a list of classifications detected in the recording.
	if len(message.Classifications) > 0 {
		body = strings.ReplaceAll(body, "{{classifications}}", strings.Join(message.Classifications, ", "))
	}

	// {{date}} {{time}} {{datetime}} of the start.
	if len(message.Media) > 0 && message.Media[0].StartTimestamp > 0 {
		t := time.Unix(message.Media[0].StartTimestamp, 0)
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
}*/

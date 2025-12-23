package channels

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	message "github.com/uug-ai/hub-pipeline-notification/message"
	templates "github.com/uug-ai/hub-pipeline-notification/templates"
	"gopkg.in/gomail.v2"
)

type SMTP struct {
	Server     string `json:"server,omitempty"`
	Port       string `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	EmailFrom  string `json:"email_from,omitempty"`
	EmailTo    string `json:"email_to,omitempty"`
	TemplateId string `json:"template_id,omitempty"`
}

func (smtp SMTP) Send(message message.Message) error {
	port, _ := strconv.Atoi(smtp.Port)
	d := gomail.NewDialer(smtp.Server, port, smtp.Username, smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", smtp.EmailFrom)
	m.SetHeader("To", smtp.EmailTo)
	m.SetHeader("Subject", message.Title)
	timeNow := time.Now().Unix()
	m.SetHeader("Message-Id", "<"+strconv.FormatInt(timeNow, 10)+"@kerberos.io>")

	textBody := templates.GetTextTemplate(smtp.TemplateId)
	m.SetBody("text/plain", ReplaceValues(textBody, message))

	body := templates.GetTemplate(smtp.TemplateId)
	m.AddAlternative("text/html", ReplaceValues(body, message))

	err := d.DialAndSend(m)
	fmt.Println(err)
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

func ReplaceValues(body string, message message.Message) string {

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
		longUrl := message.Media[0].Url
		body = strings.ReplaceAll(body, "{{link}}", longUrl)
	}
	if message.Data["link"] != "" {
		body = strings.ReplaceAll(body, "{{link}}", message.Data["link"])
	}

	// {{thumbnail}} this will inject an image (either a base64 or a url).
	if message.Thumbnail != "" {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='base64:"+message.Thumbnail+"' width='400px' height='auto' />")
	} else if len(message.Media) > 0 && message.Media[0].ThumbnailUrl != "" {
		body = strings.ReplaceAll(body, "{{thumbnail}}", "<img src='"+message.Media[0].ThumbnailUrl+"' width='400px' height='auto' />")
	}

	// {{classifications}} this will inject a list of classifications detected in the recording.
	if len(message.Classifications) > 0 {
		body = strings.ReplaceAll(body, "{{classifications}}", strings.Join(message.Classifications, ", "))
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

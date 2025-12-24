package integrations

import (
	"os"
	"testing"
	"time"

	"github.com/uug-ai/models/pkg/models"
)

var smtpTimeout = 0
var smtpTimeoutIncrement = 1000

func setupSMTPTest(t *testing.T) SMTP {
	// Initialize SMTP configuration from environment variables
	smtpMailtrap := SMTP{
		Server:    os.Getenv("SMTP_SERVER"),
		Port:      os.Getenv("SMTP_PORT"),
		Username:  os.Getenv("SMTP_USERNAME"),
		Password:  os.Getenv("SMTP_PASSWORD"),
		EmailFrom: os.Getenv("EMAIL_FROM"),
		EmailTo:   os.Getenv("EMAIL_TO"),
	}

	// Timeout, to avoid hitting issues with SMTP servers.
	smtpTimeout = smtpTimeout + smtpTimeoutIncrement
	tout := time.Duration(smtpTimeout) * time.Millisecond
	time.Sleep(tout)

	return smtpMailtrap
}

func TestSMTPWelcome(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Welcome to Kerberos Hub - Activate your account"
	m.User = "cedricve"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://kerberos.io"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "welcome"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPAssignTask(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) You have been assigned a task on Kerberos Hub"

	dataFields := make(map[string]string)
	dataFields["user"] = "kilian"
	dataFields["assignee"] = "cedricve"
	dataFields["task_name"] = "task_name"
	dataFields["link"] = "https://kerberos.io"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "assign_task"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//	t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPForgot(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Password reset Kerberos Hub. You forgot your password"
	m.User = "cedricve"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["password"] = "shizzle12345678"
	dataFields["ipaddress"] = "129.45.1.5"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "forgot"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPActivate(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Wonderful! Your Kerberos Hub is now active"
	m.User = "cedricve"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "activate"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPDetection(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Alert: Kerberos Hub detected something."

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://kerberos.io"
	dataFields["text"] = "Two dogs and a pedestrian were detected."
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "detection"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPHighUpload(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Alert: High upload triggered to Kerberos Hub"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://app.kerberos.io"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "highupload"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPDevice(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Alert: A Kerberos Agent's status changed"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://app.kerberos.io"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "device"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPDisabled(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) Daily upload reached: Kerberos Hub account disabled"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://app.kerberos.io"
	dataFields["dataUsage"] = "2"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "disable"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

func TestSMTPNewIP(t *testing.T) {
	smtpMailtrap := setupSMTPTest(t)

	m := models.Message{}
	m.Type = "message"
	m.Timestamp = int64(time.Now().Unix())
	m.Title = "(SMTP) New Login: A new device/location connect to your Kerberos Hub account"

	dataFields := make(map[string]string)
	dataFields["user"] = "cedricve"
	dataFields["link"] = "https://app.kerberos.io/profile"
	dataFields["dataUsage"] = "2"
	m.Data = dataFields

	// Send message to mailtrap.
	smtpMailtrap.TemplateId = "newip"
	err := smtpMailtrap.Send(m, "", "")
	if err != nil {
		//t.Errorf("expected error to be nil got %v", err)
	}
}

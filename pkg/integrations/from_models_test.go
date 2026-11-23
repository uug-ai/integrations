package integrations

import (
	"errors"
	"testing"

	"github.com/uug-ai/models/pkg/models"
)

func TestBuildersMapSavedChannelFields(t *testing.T) {
	if got, want := SlackFromModel(models.Slack{
		DeliveryMode: models.SlackDeliveryModeBot, Apikey: "xoxb", Channel: "C1", Hook: "https://hook", Username: "bot",
	}), (Slack{DeliveryMode: models.SlackDeliveryModeBot, Token: "xoxb", Channel: "C1", Hook: "https://hook", Username: "bot"}); got != want {
		t.Errorf("SlackFromModel = %+v, want %+v", got, want)
	}
	if got, want := WebhookFromModel(models.Webhook{Url: "https://example.com"}), (Webhook{Url: "https://example.com"}); got != want {
		t.Errorf("WebhookFromModel = %+v, want %+v", got, want)
	}
	if got, want := TelegramFromModel(models.Telegram{Token: "t", Channel: "c"}), (Telegram{Token: "t", Channel: "c"}); got != want {
		t.Errorf("TelegramFromModel = %+v, want %+v", got, want)
	}
	if got, want := PushbulletFromModel(models.Pushbullet{Apikey: "k"}), (Pushbullet{ApiKey: "k"}); got != want {
		t.Errorf("PushbulletFromModel = %+v, want %+v", got, want)
	}
	if got, want := PushoverFromModel(models.Pushover{Apikey: "k", Sendto: "u"}), (Pushover{ApiKey: "k", SendTo: "u"}); got != want {
		t.Errorf("PushoverFromModel = %+v, want %+v", got, want)
	}
	if got, want := AlexaFromModel(models.Alexa{Accesscode: "a"}), (Alexa{AccessCode: "a"}); got != want {
		t.Errorf("AlexaFromModel = %+v, want %+v", got, want)
	}
	if got, want := IftttFromModel(models.Ifttt{Token: "t"}), (Ifttt{Token: "t"}); got != want {
		t.Errorf("IftttFromModel = %+v, want %+v", got, want)
	}
	if got, want := SmsFromModel(models.Sms{Accountsid: "AC", AuthToken: "tok", From: "+1", To: "+2"}),
		(Sms{AccountSID: "AC", AuthToken: "tok", From: "+1", To: "+2"}); got != want {
		t.Errorf("SmsFromModel = %+v, want %+v", got, want)
	}
}

func TestEmailFromEnvPicksProvider(t *testing.T) {
	t.Setenv("SMTP_SERVER", "smtp.example.com")
	t.Setenv("SMTP_PORT", "2525")
	t.Setenv("SMTP_USERNAME", "user")
	t.Setenv("SMTP_PASSWORD", "pass")
	t.Setenv("EMAIL_FROM", "hub@example.com")
	t.Setenv("MAILGUN_DOMAIN", "mg.example.com")
	t.Setenv("MAILGUN_API_KEY", "mg-key")
	t.Setenv("SENDGRID_API_KEY", "sg-key")
	t.Setenv("EMAIL_FROM_DISPLAYNAME", "Hub")

	email, err := EmailFromEnv(EmailProviderSMTP, "to@example.com", "detection")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := email, (SMTP{Server: "smtp.example.com", Port: "2525", Username: "user", Password: "pass",
		EmailFrom: "hub@example.com", EmailTo: "to@example.com", TemplateId: "detection"}); got != want {
		t.Errorf("smtp = %+v, want %+v", got, want)
	}

	email, err = EmailFromEnv(EmailProviderMailgun, "to@example.com", "detection")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := email, (Mail{Domain: "mg.example.com", ApiKey: "mg-key", EmailFrom: "hub@example.com",
		EmailTo: "to@example.com", TemplateId: "detection"}); got != want {
		t.Errorf("mailgun = %+v, want %+v", got, want)
	}

	email, err = EmailFromEnv(EmailProviderSendgrid, "to@example.com", "detection")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := email, (sendgridEmail{Sendgrid{ApiKey: "sg-key", EmailFrom: "hub@example.com", EmailFromDisplay: "Hub",
		EmailTo: "to@example.com", TemplateId: "detection"}}); got != want {
		t.Errorf("sendgrid = %+v, want %+v", got, want)
	}

	_, err = EmailFromEnv("postmark", "to@example.com", "detection")
	var unsupported UnsupportedEmailProviderError
	if !errors.As(err, &unsupported) || string(unsupported) != "postmark" {
		t.Errorf("unknown provider error = %v", err)
	}
}

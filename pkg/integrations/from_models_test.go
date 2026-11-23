package integrations

import (
	"errors"
	"testing"

	"github.com/uug-ai/models/pkg/models"
)

func TestBuildersMapSavedChannelFields(t *testing.T) {
	slack, err := SlackFromModel(models.Slack{
		DeliveryMode: models.SlackDeliveryModeBot, Apikey: "xoxb", Channel: "C1", Hook: "https://hooks.example.com/1", Username: "bot",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *slack.options, (SlackOptions{DeliveryMode: models.SlackDeliveryModeBot, Token: "xoxb", Channel: "C1", Hook: "https://hooks.example.com/1", Username: "bot"}); got != want {
		t.Errorf("SlackFromModel = %+v, want %+v", got, want)
	}

	webhook, err := WebhookFromModel(models.Webhook{Url: "https://example.com"})
	if err != nil || *webhook.options != (WebhookOptions{Url: "https://example.com"}) {
		t.Errorf("WebhookFromModel = %+v, %v", webhook, err)
	}
	telegram, err := TelegramFromModel(models.Telegram{Token: "t", Channel: "c"})
	if err != nil || *telegram.options != (TelegramOptions{Token: "t", Channel: "c"}) {
		t.Errorf("TelegramFromModel = %+v, %v", telegram, err)
	}
	pushbullet, err := PushbulletFromModel(models.Pushbullet{Apikey: "k"})
	if err != nil || *pushbullet.options != (PushbulletOptions{ApiKey: "k"}) {
		t.Errorf("PushbulletFromModel = %+v, %v", pushbullet, err)
	}
	pushover, err := PushoverFromModel(models.Pushover{Apikey: "k", Sendto: "u"})
	if err != nil || *pushover.options != (PushoverOptions{ApiKey: "k", SendTo: "u"}) {
		t.Errorf("PushoverFromModel = %+v, %v", pushover, err)
	}
	alexa, err := AlexaFromModel(models.Alexa{Accesscode: "a"})
	if err != nil || *alexa.options != (AlexaOptions{AccessCode: "a"}) {
		t.Errorf("AlexaFromModel = %+v, %v", alexa, err)
	}
	ifttt, err := IftttFromModel(models.Ifttt{Token: "t"})
	if err != nil || *ifttt.options != (IftttOptions{Token: "t"}) {
		t.Errorf("IftttFromModel = %+v, %v", ifttt, err)
	}
	sms, err := SmsFromModel(models.Sms{Accountsid: "AC", AuthToken: "tok", From: "+1", To: "+2"})
	if err != nil || *sms.options != (SmsOptions{AccountSID: "AC", AuthToken: "tok", From: "+1", To: "+2"}) {
		t.Errorf("SmsFromModel = %+v, %v", sms, err)
	}
}

// A saved channel the integration cannot deliver with fails when it is built,
// as a rejected delivery, so callers record it once instead of retrying.
func TestBuildersRejectUnusableChannels(t *testing.T) {
	builds := map[string]func() error{
		"webhook without url":  func() error { _, err := WebhookFromModel(models.Webhook{}); return err },
		"webhook bad url":      func() error { _, err := WebhookFromModel(models.Webhook{Url: "not a url"}); return err },
		"slack unknown mode":   func() error { _, err := SlackFromModel(models.Slack{DeliveryMode: "unknown"}); return err },
		"telegram no channel":  func() error { _, err := TelegramFromModel(models.Telegram{Token: "t"}); return err },
		"pushover no receiver": func() error { _, err := PushoverFromModel(models.Pushover{Apikey: "k"}); return err },
		"sms no number":        func() error { _, err := SmsFromModel(models.Sms{Accountsid: "AC", AuthToken: "tok"}); return err },
	}
	for name, build := range builds {
		err := build()
		var options OptionsError
		if !errors.As(err, &options) || !options.DeliveryRejected() {
			t.Errorf("%s: error = %v, want OptionsError", name, err)
		}
	}
}

// Slack setup problems the options allow (a bot without a channel, a webhook
// mode without a hook) are still reported when sending, as before.
func TestSlackFromModelKeepsSendTimeChecks(t *testing.T) {
	if _, err := SlackFromModel(models.Slack{Apikey: "xoxb"}); err != nil {
		t.Fatalf("SlackFromModel = %v, want the send to report the missing channel", err)
	}
}

func TestNewRejectsNilOptions(t *testing.T) {
	if _, err := NewSlack(nil); err == nil {
		t.Fatal("NewSlack(nil) succeeded")
	}
}

func TestNewCopiesOptions(t *testing.T) {
	opts := NewWebhookOptions().SetUrl("https://example.com/a").Build()
	webhook, err := NewWebhook(opts)
	if err != nil {
		t.Fatal(err)
	}
	opts.Url = "https://example.com/b"
	if webhook.options.Url != "https://example.com/a" {
		t.Fatal("changing the options after NewWebhook changed the sender")
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

	email, err := EmailFromEnv(EmailProviderSMTP, "to@example.com", "detection", Templates{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *email.(*SMTP).options, (SMTPOptions{Server: "smtp.example.com", Port: 2525, Username: "user", Password: "pass",
		EmailFrom: "hub@example.com", EmailTo: "to@example.com", TemplateId: "detection"}); got != want {
		t.Errorf("smtp = %+v, want %+v", got, want)
	}

	email, err = EmailFromEnv(EmailProviderMailgun, "to@example.com", "detection", Templates{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *email.(*Mail).options, (MailOptions{Domain: "mg.example.com", ApiKey: "mg-key", EmailFrom: "hub@example.com",
		EmailTo: "to@example.com", TemplateId: "detection"}); got != want {
		t.Errorf("mailgun = %+v, want %+v", got, want)
	}

	email, err = EmailFromEnv(EmailProviderSendgrid, "to@example.com", "detection", Templates{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := *email.(sendgridEmail).sendgrid.options, (SendgridOptions{ApiKey: "sg-key", EmailFrom: "hub@example.com", EmailFromDisplay: "Hub",
		EmailTo: "to@example.com", TemplateId: "detection"}); got != want {
		t.Errorf("sendgrid = %+v, want %+v", got, want)
	}

	_, err = EmailFromEnv("postmark", "to@example.com", "detection", Templates{})
	var unsupported UnsupportedEmailProviderError
	if !errors.As(err, &unsupported) || string(unsupported) != "postmark" {
		t.Errorf("unknown provider error = %v", err)
	}
}

func TestEmailFromEnvWithoutSettingsIsRejected(t *testing.T) {
	for _, name := range []string{"SMTP_SERVER", "SMTP_PORT", "EMAIL_FROM"} {
		t.Setenv(name, "")
	}
	email, err := EmailFromEnv(EmailProviderSMTP, "to@example.com", "detection", Templates{})
	var options OptionsError
	if !errors.As(err, &options) {
		t.Fatalf("error = %v, want OptionsError", err)
	}
	if email != nil {
		t.Fatalf("email = %#v, want nil with the error", email)
	}
}

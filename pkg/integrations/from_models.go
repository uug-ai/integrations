package integrations

import "github.com/uug-ai/models/pkg/models"

// The builders below turn a saved channel (models.Channels) into the sender
// that delivers to it, so every service maps the stored fields the same way.
// They do not check Enabled or Valid; callers decide which channels to send
// to. Options the integration cannot deliver with return an OptionsError.

// SlackFromModel builds the Slack sender for a saved Slack channel. The saved
// delivery mode picks between the bot and the webhook; an empty mode falls
// back to the fields that are set.
func SlackFromModel(slack models.Slack) (*Slack, error) {
	return NewSlack(NewSlackOptions().
		SetDeliveryMode(slack.DeliveryMode).
		SetToken(slack.Apikey).
		SetChannel(slack.Channel).
		SetHook(slack.Hook).
		SetUsername(slack.Username).
		Build())
}

// WebhookFromModel builds the webhook sender for a saved webhook channel.
func WebhookFromModel(webhook models.Webhook) (*Webhook, error) {
	return NewWebhook(NewWebhookOptions().SetUrl(webhook.Url).Build())
}

// TelegramFromModel builds the Telegram sender for a saved Telegram channel.
func TelegramFromModel(telegram models.Telegram) (*Telegram, error) {
	return NewTelegram(NewTelegramOptions().
		SetToken(telegram.Token).
		SetChannel(telegram.Channel).
		Build())
}

// PushbulletFromModel builds the Pushbullet sender for a saved Pushbullet channel.
func PushbulletFromModel(pushbullet models.Pushbullet) (*Pushbullet, error) {
	return NewPushbullet(NewPushbulletOptions().SetApiKey(pushbullet.Apikey).Build())
}

// PushoverFromModel builds the Pushover sender for a saved Pushover channel.
func PushoverFromModel(pushover models.Pushover) (*Pushover, error) {
	return NewPushover(NewPushoverOptions().
		SetApiKey(pushover.Apikey).
		SetSendTo(pushover.Sendto).
		Build())
}

// AlexaFromModel builds the Alexa sender for a saved Alexa channel.
func AlexaFromModel(alexa models.Alexa) (*Alexa, error) {
	return NewAlexa(NewAlexaOptions().SetAccessCode(alexa.Accesscode).Build())
}

// IftttFromModel builds the IFTTT sender for a saved IFTTT channel.
func IftttFromModel(ifttt models.Ifttt) (*Ifttt, error) {
	return NewIfttt(NewIftttOptions().SetToken(ifttt.Token).Build())
}

// SmsFromModel builds the Twilio SMS sender for a saved SMS channel.
func SmsFromModel(sms models.Sms) (*Sms, error) {
	return NewSms(NewSmsOptions().
		SetAccountSID(sms.Accountsid).
		SetAuthToken(sms.AuthToken).
		SetFrom(sms.From).
		SetTo(sms.To).
		Build())
}

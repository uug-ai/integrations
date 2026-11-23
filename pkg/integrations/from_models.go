package integrations

import "github.com/uug-ai/models/pkg/models"

// The builders below turn a saved channel (models.Channels) into the sender
// that delivers to it, so every service maps the stored fields the same way.
// They do not check Enabled or Valid; callers decide which channels to send to.

// SlackFromModel builds the Slack sender for a saved Slack channel. The saved
// delivery mode picks between the bot and the webhook; an empty mode falls
// back to the fields that are set.
func SlackFromModel(slack models.Slack) Slack {
	return Slack{
		DeliveryMode: slack.DeliveryMode,
		Token:        slack.Apikey,
		Channel:      slack.Channel,
		Hook:         slack.Hook,
		Username:     slack.Username,
	}
}

// WebhookFromModel builds the webhook sender for a saved webhook channel.
func WebhookFromModel(webhook models.Webhook) Webhook {
	return Webhook{Url: webhook.Url}
}

// TelegramFromModel builds the Telegram sender for a saved Telegram channel.
func TelegramFromModel(telegram models.Telegram) Telegram {
	return Telegram{Token: telegram.Token, Channel: telegram.Channel}
}

// PushbulletFromModel builds the Pushbullet sender for a saved Pushbullet channel.
func PushbulletFromModel(pushbullet models.Pushbullet) Pushbullet {
	return Pushbullet{ApiKey: pushbullet.Apikey}
}

// PushoverFromModel builds the Pushover sender for a saved Pushover channel.
func PushoverFromModel(pushover models.Pushover) Pushover {
	return Pushover{ApiKey: pushover.Apikey, SendTo: pushover.Sendto}
}

// AlexaFromModel builds the Alexa sender for a saved Alexa channel.
func AlexaFromModel(alexa models.Alexa) Alexa {
	return Alexa{AccessCode: alexa.Accesscode}
}

// IftttFromModel builds the IFTTT sender for a saved IFTTT channel.
func IftttFromModel(ifttt models.Ifttt) Ifttt {
	return Ifttt{Token: ifttt.Token}
}

// SmsFromModel builds the Twilio SMS sender for a saved SMS channel.
func SmsFromModel(sms models.Sms) Sms {
	return Sms{
		AccountSID: sms.Accountsid,
		AuthToken:  sms.AuthToken,
		From:       sms.From,
		To:         sms.To,
	}
}

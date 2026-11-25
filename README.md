# Integrations

Notification channels for Go: Slack, webhooks, email (SMTP, Mailgun, SendGrid), Telegram, Pushover, Pushbullet, Pusher, Alexa, IFTTT, SMS (Twilio) and MQTT.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.26-007d9c?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/uug-ai/integrations)](https://goreportcard.com/report/github.com/uug-ai/integrations)
[![codecov](https://codecov.io/gh/uug-ai/integrations/graph/badge.svg?token=BOInBI4j2N)](https://codecov.io/gh/uug-ai/integrations)

## Installation

```bash
go get github.com/uug-ai/integrations
```

## Packages

| Package | Contents |
|---------|----------|
| `pkg/message` | `Message`, the payload every channel sends, and `Media`. |
| `pkg/integrations` | One sender per channel, builders from saved channels (`models.Channels`), email providers configured from the environment, and the SMTP template lookup. |

## Sending a message

Every integration is created the same way: build its options with `New<Integration>Options()`, chaining `Set<Field>` calls and ending with `Build()`, then pass them to `New<Integration>(opts)`. The constructor validates the options and returns the sender, or an `OptionsError` when the integration cannot deliver with them (a missing token, an invalid URL). The sender keeps a copy of the options.

```go
import (
    "context"

    "github.com/uug-ai/integrations/pkg/integrations"
    "github.com/uug-ai/integrations/pkg/message"
)

m := message.Message{
    Title: "Motion detected",
    Body:  "Front door, 14:53",
    Media: []message.Media{{Type: "video", Url: videoURL, ThumbnailUrl: thumbnailURL}},
}

slack, err := integrations.NewSlack(integrations.NewSlackOptions().
    SetHook("https://hooks.slack.com/services/...").
    SetUsername("Hub").
    Build())
if err != nil {
    return err
}
err = slack.SendWithContext(ctx, m)

webhook, err := integrations.NewWebhook(integrations.NewWebhookOptions().
    SetUrl("https://example.com/hook").
    Build())
if err != nil {
    return err
}
err = webhook.SendWithContext(ctx, m)
```

| Integration | Options builder | Constructor | Send |
|-------------|-----------------|-------------|------|
| Slack | `NewSlackOptions()` | `NewSlack` | `Send`, `SendWithContext` |
| Webhook | `NewWebhookOptions()` | `NewWebhook` | `Send`, `SendWithContext` |
| SMTP | `NewSMTPOptions()` | `NewSMTP` | `Send` |
| Mailgun | `NewMailOptions()` | `NewMail` | `Send` |
| SendGrid | `NewSendgridOptions()` | `NewSendgrid` | `SendNotification` |
| Telegram | `NewTelegramOptions()` | `NewTelegram` | `Send` |
| Pushover | `NewPushoverOptions()` | `NewPushover` | `Send` |
| Pushbullet | `NewPushbulletOptions()` | `NewPushbullet` | `SendLink`, `SendMessage` |
| Pusher | `NewPusherOptions()` | `NewPusher` | `Trigger`, `Send`, `SendNotification` |
| Alexa | `NewAlexaOptions()` | `NewAlexa` | `Send` |
| IFTTT | `NewIftttOptions()` | `NewIfttt` | `Send` |
| SMS (Twilio) | `NewSmsOptions()` | `NewSms` | `Send` |
| MQTT | `NewMQTTOptions()` | `NewMQTT` | `Publish`, `SendNotification` |

`OptionsError` reports `DeliveryRejected() true`, so a caller that classifies delivery errors treats bad options as a rejection, not something to retry.

### From saved channels

The `...FromModel` helpers build the sender for a saved channel from `github.com/uug-ai/models` through its options builder. They do not check `Enabled` or `Valid`; the caller picks the channels to send to.

```go
if channels.Slack.Enabled {
    slack, err := integrations.SlackFromModel(channels.Slack)
    if err == nil {
        err = slack.SendWithContext(ctx, m)
    }
}
```

Helpers: `SlackFromModel`, `WebhookFromModel`, `TelegramFromModel`, `PushbulletFromModel`, `PushoverFromModel`, `AlexaFromModel`, `IftttFromModel`, `SmsFromModel`.

### Email

`EmailFromEnv(provider, to, template, templates)` returns the sender for `smtp`, `mailgun` or `sendgrid`, configured from the environment through the options builders. `SMTPFromEnv`, `MailgunFromEnv` and `SendgridFromEnv` build one provider directly. Missing settings (server, port, sender, API key) return an `OptionsError`.

| Provider | Environment |
|----------|-------------|
| smtp | `SMTP_SERVER`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `EMAIL_FROM` |
| mailgun | `MAILGUN_DOMAIN`, `MAILGUN_API_KEY`, `EMAIL_FROM` |
| sendgrid | `SENDGRID_API_KEY`, `EMAIL_FROM`, `EMAIL_FROM_DISPLAYNAME` |

SMTP renders the template named by `TemplateId` from its `Templates`. Variables such as `{{text}}`, `{{link}}`, `{{thumbnail}}`, `{{devicename}}` and every key of `Message.Data` are filled in (see `ReplaceValues`). Mailgun and SendGrid pass the same values to a template stored at the provider.

The templates belong to the service that sends them. A service embeds its own `<name>.html` and `<name>.txt` files, and `NewTemplates` lets files in `/mail/templates` (`DefaultTemplateDir`) replace them one by one:

```go
//go:embed templates/*.html templates/*.txt
var files embed.FS

builtin, _ := fs.Sub(files, "templates")
templates := integrations.NewTemplates(builtin)

email, err := integrations.EmailFromEnv(os.Getenv("MAIL_PROVIDER"), "user@example.com", "detection", templates)
if err == nil {
    err = email.Send(m)
}
```

### Slack delivery modes

`DeliveryMode` picks how Slack is reached:

- `webhook` posts to `Hook`.
- `bot` posts with the bot `Token` to `Channel`. A message whose first media item is a video uploads the video file (up to 1 GB) to the channel.

An empty mode falls back to the fields that are set.

### Delivery errors

Errors keep enough detail to decide whether a retry can succeed:

- Webhook errors for a non-2xx answer implement `HTTPStatusCode() int`.
- Slack setup errors (missing token, channel or hook, a video that is empty or too large) implement `DeliveryRejected() bool`.
- Slack API and SMTP errors are returned as their libraries report them.

### Pusher

`Pusher` needs the Pusher app credentials (`SetAppID`, `SetKey`, `SetSecret`, `SetCluster`); `PusherFromEnv(event)` reads them from `PUSHER_APP_ID`, `PUSHER_KEY`, `PUSHER_SECRET` and `PUSHER_CLUSTER`. It triggers the `Channel` event on the channel named after `Message.User`, which must be the recipient's username. `Trigger` returns the delivery error; `Send` only reports whether the recipient was valid.

### MQTT

`MQTT` needs the broker URI (`SetURI`, for example `tcp://mqtt.kerberos.io:1883`) and, when the broker asks for them, `SetUsername` and `SetPassword`. `Publish(topic, payload)` connects, publishes the payload on the topic with QoS 1 (not retained), waits up to 10 seconds for the broker to acknowledge it, and disconnects; it returns an error when the broker cannot be reached or does not acknowledge in time. `SendNotification` connects to the broker but does not publish the message.

## Testing

```bash
go test ./...
```

Tests use local HTTP servers and do not reach the providers.

## Contributing

1. Create a branch (`feat/...`, `fix/...`).
2. Give a new integration the same shape as the others: an `<Name>Options` struct with `validate` tags, a `<Name>OptionsBuilder` from `New<Name>Options()` with one `Set<Field>` per field and `Build()`, and `New<Name>(opts)` that calls `validateOptions` and keeps a copy of the options.
3. Add tests next to the code.
4. Run `go vet ./...` and `go test ./...`.
5. Write commit messages following [Conventional Commits](https://www.conventionalcommits.org/), with the integration as scope when it applies (`fix(slack): ...`).

## License

MIT, see [LICENSE](LICENSE).

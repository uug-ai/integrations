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
| `pkg/integrations` | One sender per channel, builders from saved channels (`models.Channels`), and email providers configured from the environment. |
| `pkg/templates` | Built-in email templates (HTML and plain text). A file in `/mail/templates/<name>.html` or `<name>.txt` replaces the built-in one. |

## Sending a message

Each channel is a plain struct with its settings and a `Send` method that takes a `message.Message`.

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

slack := integrations.Slack{Hook: "https://hooks.slack.com/services/...", Username: "Hub"}
err := slack.SendWithContext(ctx, m)

webhook := integrations.Webhook{Url: "https://example.com/hook"}
err = webhook.SendWithContext(ctx, m)
```

### From saved channels

The `...FromModel` builders map a saved channel from `github.com/uug-ai/models` to its sender. They do not check `Enabled` or `Valid`; the caller picks the channels to send to.

```go
if channels.Slack.Enabled {
    err := integrations.SlackFromModel(channels.Slack).SendWithContext(ctx, m)
}
```

Builders: `SlackFromModel`, `WebhookFromModel`, `TelegramFromModel`, `PushbulletFromModel`, `PushoverFromModel`, `AlexaFromModel`, `IftttFromModel`, `SmsFromModel`.

### Email

`EmailFromEnv(provider, to, template)` returns the sender for `smtp`, `mailgun` or `sendgrid`, configured from the environment. `SMTPFromEnv`, `MailgunFromEnv` and `SendgridFromEnv` build one provider directly.

| Provider | Environment |
|----------|-------------|
| smtp | `SMTP_SERVER`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `EMAIL_FROM` |
| mailgun | `MAILGUN_DOMAIN`, `MAILGUN_API_KEY`, `EMAIL_FROM` |
| sendgrid | `SENDGRID_API_KEY`, `EMAIL_FROM`, `EMAIL_FROM_DISPLAYNAME` |

```go
email, err := integrations.EmailFromEnv(os.Getenv("MAIL_PROVIDER"), "user@example.com", "detection")
if err == nil {
    err = email.Send(m)
}
```

SMTP renders the template named by `TemplateId` from `pkg/templates`, filling variables such as `{{text}}`, `{{link}}`, `{{thumbnail}}`, `{{devicename}}` and every key of `Message.Data` (see `ReplaceValues`). Mailgun and SendGrid pass the same values to a template stored at the provider.

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

`Pusher` needs the Pusher app credentials (`AppID`, `Key`, `Secret`, `Cluster`); `PusherFromEnv(event)` reads them from `PUSHER_APP_ID`, `PUSHER_KEY`, `PUSHER_SECRET` and `PUSHER_CLUSTER`. It triggers the `Channel` event on the channel named after `Message.User`, which must be the recipient's username. `Trigger` returns the delivery error; `Send` only reports whether the recipient was valid.

## Testing

```bash
go test ./...
```

Tests use local HTTP servers and do not reach the providers.

## Contributing

1. Create a branch (`feat/...`, `fix/...`).
2. Add tests next to the code.
3. Run `go vet ./...` and `go test ./...`.
4. Write commit messages following [Conventional Commits](https://www.conventionalcommits.org/), with the integration as scope when it applies (`fix(slack): ...`).

## License

MIT, see [LICENSE](LICENSE).

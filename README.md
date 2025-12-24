# Integrations

Universal notification and messaging integrations for Go.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-007d9c?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/uug-ai/integrations)](https://goreportcard.com/report/github.com/uug-ai/integrations)
[![codecov](https://codecov.io/gh/uug-ai/integrations/graph/badge.svg?token=BOInBI4j2N)](https://codecov.io/gh/uug-ai/integrations)

A Go library for sending notifications and messages across multiple platforms and services with a unified interface.

## Features

- 15+ integrations for popular messaging and notification services
- Unified interface across all integrations
- Type safe with full Go type safety and compile-time checks
- Comprehensive test coverage
- MongoDB support
- OpenTelemetry observability and tracing support

## Installation

```bash
go get github.com/uug-ai/integrations
```

## Quick Start

```go
package main

import (
    "github.com/uug-ai/integrations/pkg/integrations"
)

func main() {
    // Slack example
    slack := integrations.Slack{
        Hook:     "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
        Username: "MyBot",
    }

    body := "Hello from integrations!"
    imageUrl := "https://example.com/image.png" // Optional image URL, use empty string if not needed

    err := slack.Send(body, imageUrl)
    if err != nil {
        panic(err)
    }
}
```

## Supported Integrations

| Integration | Description | Status |
|-------------|-------------|--------|
| **Slack** | Send messages to Slack channels | Ready |
| **Telegram** | Send messages via Telegram Bot API | Ready |
| **Webhook** | Generic HTTP webhooks | Ready |
| **SMTP** | Email via SMTP | Ready |
| **SendGrid** | Email via SendGrid API | Ready |
| **Pushover** | Push notifications to mobile devices | Ready |
| **Pushbullet** | Push notifications and file sharing | Ready |
| **Pusher** | Real-time messaging via Pusher | Ready |
| **MQTT** | IoT messaging protocol | Ready |
| **Alexa** | Amazon Alexa notifications | Ready |
| **IFTTT** | If This Then That automation | Ready |
| **Twitter** | Tweet updates (via Twitter API) | Ready |
| **SMS** | SMS via Twilio | Ready |
| **Mail** | Email via Mailgun | Ready |
| **MongoDB** | Database storage integration | Ready |

## Usage Examples

### Slack

```go
slack := integrations.Slack{
    Hook:     "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
    Username: "MyBot",
}

// Send a text message
err := slack.Send("Hello from Slack!", "")

// Send a message with an image attachment
err := slack.Send("Check out this image!", "https://example.com/image.png")
```

## Project Structure

```
.
├── pkg/
│   └── integrations/        # Core integration implementations
│       ├── alexa.go
│       ├── ifttt.go
│       ├── mail.go
│       ├── message.go       # Message struct definition
│       ├── mongodb.go
│       ├── mqtt.go
│       ├── pushbullet.go
│       ├── pusher.go
│       ├── pushover.go
│       ├── sendgrid.go
│       ├── slack.go
│       ├── sms.go
│       ├── smtp.go
│       ├── telegram.go
│       ├── twitter.go
│       └── webhook.go
├── main.go
├── go.mod
└── README.md
```

## Slack Configuration

The Slack integration uses webhook URLs for sending messages:

```go
type Slack struct {
    Hook     string // Slack webhook URL
    Username string // Bot username to display
}
```

The `Send` method accepts two parameters:
- `body` (string): The message text to send
- `url` (string): Optional image URL to attach (use empty string if not needed)

## Configuration

Each integration has its own configuration struct. Check the individual integration files in `pkg/integrations/` for specific configuration options.

### Environment Variables

Many integrations support configuration via environment variables:

```bash
# Slack
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."

# Telegram
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHANNEL_ID="your_channel"

# Pushover
export PUSHOVER_TOKEN="your_app_token"
export PUSHOVER_USER="your_user_key"
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests for a specific integration:

```bash
go test ./pkg/integrations -run TestSlack
```

## Contributing

Contributions are welcome. To contribute:

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/)
4. Push to your branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`

**Examples:**
```
feat(telegram): add support for inline keyboards
fix(slack): correct webhook payload formatting
docs(readme): update usage examples
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project uses the following libraries:

- [slack-go/slack](https://github.com/slack-go/slack) - Slack API in Go
- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API
- [sendgrid-go](https://github.com/sendgrid/sendgrid-go) - SendGrid email API
- [mongo-driver](https://github.com/mongodb/mongo-go-driver) - MongoDB driver
- [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) - MQTT client

See [go.mod](go.mod) for the complete list of dependencies.

## Support

- Issues: [GitHub Issues](https://github.com/uug-ai/integrations/issues)
- Discussions: [GitHub Discussions](https://github.com/uug-ai/integrations/discussions)
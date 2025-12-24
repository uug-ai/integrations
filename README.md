# Integrations

Universal notification and messaging integrations for Go.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-007d9c?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/uug-ai/integrations)](https://goreportcard.com/report/github.com/uug-ai/integrations)
[![codecov](https://codecov.io/gh/uug-ai/integrations/graph/badge.svg?token=BOInBI4j2N)](https://codecov.io/gh/uug-ai/integrations)

A Go library for sending notifications and messages across multiple platforms and services with a unified interface using the functional options pattern.

## Features

- 15+ integrations for popular messaging and notification services
- Functional options pattern for flexible configuration
- Built-in validation with compile-time type safety
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
    "log"
    "github.com/uug-ai/integrations/pkg/integrations"
)

func main() {
    // Create a Slack integration using the functional options pattern
    slack, err := integrations.CreateSlack(
        integrations.WithSlackHook("https://hooks.slack.com/services/YOUR/WEBHOOK/URL"),
        integrations.WithSlackUsername("MyBot"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Send a message
    err = slack.Send("Hello from integrations!", "")
    if err != nil {
        log.Fatal(err)
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

## Core Concepts

### Functional Options Pattern

All integrations use the functional options pattern for configuration. This provides:
- **Flexibility**: Configure only what you need
- **Validation**: Built-in validation before use
- **Type Safety**: Compile-time type checking
- **Extensibility**: Easy to add new options

### Creating Integrations

Each integration follows this pattern:

1. **Create** the integration using a `Create<Integration>()` function with functional options
2. **Validate** configuration automatically during creation
3. **Send** messages using the `Send()` method

## Usage Examples

### SMTP (Email)

The SMTP integration demonstrates the functional options pattern:

```go
package main

import (
    "log"
    "github.com/uug-ai/integrations/pkg/integrations"
)

func main() {
    // Create SMTP integration with functional options
    smtp, err := integrations.CreateSMTP(
        integrations.WithSMTPServer("smtp.gmail.com"),
        integrations.WithSMTPPort(587),
        integrations.WithSMTPUsername("your-email@gmail.com"),
        integrations.WithSMTPPassword("your-app-password"),
        integrations.WithSMTPEmailFrom("sender@example.com"),
        integrations.WithSMTPEmailTo("recipient@example.com"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Send an email
    err = smtp.Send(
        "Email Subject",           // title
        "Plain text body",          // body
        "<h1>HTML body</h1>",      // textBody (HTML alternative)
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

**Available Options:**
- `WithSMTPServer(server string)` - SMTP server hostname
- `WithSMTPPort(port int)` - SMTP server port
- `WithSMTPUsername(username string)` - Authentication username
- `WithSMTPPassword(password string)` - Authentication password
- `WithSMTPEmailFrom(email string)` - Sender email address
- `WithSMTPEmailTo(email string)` - Recipient email address

### Slack

```go
// Create Slack integration
slack, err := integrations.CreateSlack(
    integrations.WithSlackHook("https://hooks.slack.com/services/YOUR/WEBHOOK/URL"),
    integrations.WithSlackUsername("MyBot"),
)
if err != nil {
    log.Fatal(err)
}

// Send a text message
err = slack.Send("Hello from Slack!", "")

// Send a message with an image attachment
err = slack.Send("Check out this image!", "https://example.com/image.png")
```

**Available Options:**
- `WithSlackHook(hook string)` - Slack webhook URL
- `WithSlackUsername(username string)` - Bot username to display

## Project Structure

```
.
├── pkg/
│   └── integrations/        # Core integration implementations
│       ├── alexa.go
│       ├── ifttt.go
│       ├── mail.go
│       ├── mongodb.go
│       ├── mqtt.go
│       ├── option.go        # Generic functional option type
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

## Creating a New Integration

To create a new integration following the functional options pattern:

### 1. Define the Integration Struct

```go
package integrations

type MyService struct {
    APIKey   string `json:"api_key" validate:"required"`
    Endpoint string `json:"endpoint" validate:"required,url"`
}
```

### 2. Create Functional Options

```go
// WithMyServiceAPIKey sets the API key
func WithMyServiceAPIKey(apiKey string) Option[MyService] {
    return func(s *MyService) {
        s.APIKey = apiKey
    }
}

// WithMyServiceEndpoint sets the endpoint URL
func WithMyServiceEndpoint(endpoint string) Option[MyService] {
    return func(s *MyService) {
        s.Endpoint = endpoint
    }
}
```

### 3. Implement the Create Function

```go
// CreateMyService creates a new MyService instance with the provided options
func CreateMyService(opts ...Option[MyService]) (*MyService, error) {
    service := &MyService{}

    // Apply all options
    for _, opt := range opts {
        opt(service)
    }

    // Validate configuration
    err := service.Validate()
    if err != nil {
        return nil, err
    }

    return service, nil
}
```

### 4. Add Validation

```go
func (s *MyService) Validate() error {
    validate := validator.New()
    err := validate.Struct(s)
    if err != nil {
        return err
    }
    return nil
}
```

### 5. Implement the Send Method

```go
func (s *MyService) Send(message string) error {
    // Implementation here
    return nil
}
```

### 6. Usage

```go
service, err := integrations.CreateMyService(
    integrations.WithMyServiceAPIKey("your-api-key"),
    integrations.WithMyServiceEndpoint("https://api.example.com"),
)
if err != nil {
    log.Fatal(err)
}

err = service.Send("Hello, World!")
```

## Validation

All integrations use [go-playground/validator](https://github.com/go-playground/validator) for configuration validation. Common validation tags:

- `required` - Field must not be empty
- `email` - Must be a valid email address
- `url` - Must be a valid URL
- `gt=0` - Must be greater than 0
- `min=<value>` - Minimum value/length
- `max=<value>` - Maximum value/length

The `Validate()` method is automatically called during the `Create<Integration>()` function, ensuring invalid configurations are caught before use.

## Configuration

### Using Functional Options (Recommended)

```go
smtp, err := integrations.CreateSMTP(
    integrations.WithSMTPServer("smtp.gmail.com"),
    integrations.WithSMTPPort(587),
    integrations.WithSMTPUsername("user@example.com"),
    integrations.WithSMTPPassword("password"),
    integrations.WithSMTPEmailFrom("from@example.com"),
    integrations.WithSMTPEmailTo("to@example.com"),
)
```
### Environment Variables

You can load configuration from environment variables before creating integrations:

```go
import "os"

smtp, err := integrations.CreateSMTP(
    integrations.WithSMTPServer(os.Getenv("SMTP_SERVER")),
    integrations.WithSMTPPort(587),
    integrations.WithSMTPUsername(os.Getenv("SMTP_USERNAME")),
    integrations.WithSMTPPassword(os.Getenv("SMTP_PASSWORD")),
    integrations.WithSMTPEmailFrom(os.Getenv("SMTP_FROM")),
    integrations.WithSMTPEmailTo(os.Getenv("SMTP_TO")),
)
```

Example `.env` file:

```bash
# SMTP Configuration
SMTP_SERVER=smtp.gmail.com
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=sender@example.com
SMTP_TO=recipient@example.com

# Slack Configuration
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
SLACK_USERNAME=MyBot

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_token
TELEGRAM_CHANNEL_ID=your_channel
```

## Error Handling

The functional options pattern provides clear error handling:

```go
smtp, err := integrations.CreateSMTP(
    integrations.WithSMTPServer("smtp.gmail.com"),
    integrations.WithSMTPPort(587),
    // Missing required fields...
)
if err != nil {
    // Validation error caught at creation time
    log.Printf("Configuration error: %v", err)
    return
}

// If we get here, the configuration is valid
err = smtp.Send("Subject", "Body", "<h1>HTML</h1>")
if err != nil {
    // Runtime error during send
    log.Printf("Send error: %v", err)
    return
}
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
# SMTP tests
go test ./pkg/integrations -run TestSMTP

# Slack tests
go test ./pkg/integrations -run TestSlack
```

Run all integration tests:

```bash
go test ./pkg/integrations/... -v
```

## Contributing

Contributions are welcome! When adding new integrations, please follow the functional options pattern demonstrated in this repository.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Follow the functional options pattern (see "Creating a New Integration" above)
4. Add comprehensive tests for your integration
5. Ensure all tests pass: `go test ./...`
6. Add usage examples to the README
7. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/)
8. Push to your branch (`git push origin feat/amazing-feature`)
9. Open a Pull Request

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `types`

**Scopes:**
- `models` - Changes to data structures
- `api` - API-related changes
- `types` - Type definitions
- `docs` - Documentation updates
- Integration names (e.g., `smtp`, `slack`, `telegram`)

**Examples:**
```
feat(smtp): add functional options pattern
feat(telegram): add CreateTelegram with options
fix(slack): correct webhook payload formatting
docs(readme): update usage examples with functional options
refactor(pushover): migrate to functional options pattern
test(smtp): add validation tests
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

This project uses the following key libraries:

- [go-playground/validator](https://github.com/go-playground/validator) - Struct validation
- [slack-go/slack](https://github.com/slack-go/slack) - Slack API in Go
- [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API
- [sendgrid-go](https://github.com/sendgrid/sendgrid-go) - SendGrid email API
- [mongo-driver](https://github.com/mongodb/mongo-go-driver) - MongoDB driver
- [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) - MQTT client
- [gomail](https://gopkg.in/gomail.v2) - SMTP email library

See [go.mod](go.mod) for the complete list of dependencies.

## Benefits of the Functional Options Pattern

### Type Safety
The generic `Option[T]` type provides compile-time type checking, preventing configuration errors.

### Flexibility
Configure only the options you need. No need to pass empty or default values.

### Validation
Built-in validation ensures configurations are correct before use, catching errors early.

### Extensibility
Adding new options doesn't break existing code. Simply add new `With*` functions.

### Readability
```go
// Clear and self-documenting
smtp, err := integrations.CreateSMTP(
    integrations.WithSMTPServer("smtp.gmail.com"),
    integrations.WithSMTPPort(587),
    integrations.WithSMTPUsername("user@example.com"),
    integrations.WithSMTPPassword("password"),
    integrations.WithSMTPEmailFrom("from@example.com"),
    integrations.WithSMTPEmailTo("to@example.com"),
)
```

## Support

- Issues: [GitHub Issues](https://github.com/uug-ai/integrations/issues)
- Discussions: [GitHub Discussions](https://github.com/uug-ai/integrations/discussions)
# Integrations

Universal notification and messaging integrations for Go.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-007d9c?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/uug-ai/integrations)](https://goreportcard.com/report/github.com/uug-ai/integrations)
[![codecov](https://codecov.io/gh/uug-ai/integrations/graph/badge.svg?token=BOInBI4j2N)](https://codecov.io/gh/uug-ai/integrations)

A Go library for sending notifications and messages across multiple platforms and services with a unified interface using the functional options pattern.

## Features

- 15+ integrations for popular messaging and notification services
- Fluent builder pattern for clean, readable configuration
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
    // Build SMTP options
    opts := integrations.NewSMTPOptions().
        Server("smtp.gmail.com").
        Port(587).
        Username("your-email@gmail.com").
        Password("your-app-password").
        From("sender@example.com").
        To("recipient@example.com").
        Build()

    // Create SMTP client with options
    smtp, err := integrations.NewSMTP(opts)
    if err != nil {
        log.Fatal(err)
    }

    // Send an email
    err = smtp.Send(
        "Email Subject",
        "Plain text body",
        "<h1>HTML body</h1>",
    )
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

### Options Builder Pattern

All integrations use the options builder pattern (similar to MongoDB's driver). This provides:
- **Clean Syntax**: Build options separately, then pass to constructor
- **Readability**: Self-documenting method chains
- **Separation of Concerns**: Options building is separate from client creation
- **Validation**: Built-in validation when creating the client
- **Type Safety**: Compile-time type checking
- **Flexibility**: Configure only what you need

### Creating Integrations

Each integration follows this pattern:

1. **Build Options** using `integrations.New<Integration>Options()` with method chaining
2. **Call** `.Build()` to get the options object
3. **Create Client** by passing options to `integrations.New<Integration>(opts)`
4. **Send** messages using the `Send()` method

## Usage Examples

### SMTP (Email)

The SMTP integration demonstrates the options builder pattern:

```go
package main

import (
    "log"
    "github.com/uug-ai/integrations/pkg/integrations"
)

func main() {
    // Build SMTP options
    opts := integrations.NewSMTPOptions().
        Server("smtp.gmail.com").
        Port(587).
        Username("your-email@gmail.com").
        Password("your-app-password").
        From("sender@example.com").
        To("recipient@example.com").
        Build()

    // Create SMTP client with options
    smtp, err := integrations.NewSMTP(opts)
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

**Available Methods:**
- `.Server(server string)` - SMTP server hostname
- `.Port(port int)` - SMTP server port
- `.Username(username string)` - Authentication username
- `.Password(password string)` - Authentication password
- `.From(email string)` - Sender email address
- `.To(email string)` - Recipient email address
- `.Build()` - Returns the SMTPOptions object

### Slack

```go
// Build Slack options
opts := integrations.NewSlackOptions().
    Hook("https://hooks.slack.com/services/YOUR/WEBHOOK/URL").
    Username("MyBot").
    Build()

// Create Slack client
slack, err := integrations.NewSlack(opts)
if err != nil {
    log.Fatal(err)
}

// Send a text message
err = slack.Send("Hello from Slack!", "")

// Send a message with an image attachment
err = slack.Send("Check out this image!", "https://example.com/image.png")
```

**Available Methods:**
- `.Hook(hook string)` - Slack webhook URL
- `.Username(username string)` - Bot username to display
- `.Build()` - Returns the SlackOptions object

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

To create a new integration following the options builder pattern:

### 1. Define the Options Struct (unexported)

```go
package integrations

// MyServiceOptions holds the configuration for MyService
type MyServiceOptions struct {
    apiKey   string `validate:"required"`
    endpoint string `validate:"required,url"`
}
```

### 2. Define the Client Struct

```go
// MyService represents a MyService client instance
type MyService struct {
    options *MyServiceOptions
}
```

### 3. Create the Options Builder Struct

```go
// MyServiceOptionsBuilder provides a fluent interface for building MyService options
type MyServiceOptionsBuilder struct {
    options *MyServiceOptions
}
```

### 4. Create the Constructor Function

```go
// NewMyServiceOptions creates a new MyService options builder
func NewMyServiceOptions() *MyServiceOptionsBuilder {
    return &MyServiceOptionsBuilder{
        options: &MyServiceOptions{},
    }
}
```

### 5. Add Builder Methods

```go
// APIKey sets the API key
func (b *MyServiceOptionsBuilder) APIKey(apiKey string) *MyServiceOptionsBuilder {
    b.options.apiKey = apiKey
    return b
}

// Endpoint sets the endpoint URL
func (b *MyServiceOptionsBuilder) Endpoint(endpoint string) *MyServiceOptionsBuilder {
    b.options.endpoint = endpoint
    return b
}

// Build returns the configured MyServiceOptions
func (b *MyServiceOptionsBuilder) Build() *MyServiceOptions {
    return b.options
}
```

### 6. Create the Client Constructor with Validation

```go
// NewMyService creates a new MyService client with the provided options
func NewMyService(opts *MyServiceOptions) (*MyService, error) {
    // Validate configuration
    validate := validator.New()
    err := validate.Struct(opts)
    if err != nil {
        return nil, err
    }

    return &MyService{
        options: opts,
    }, nil
}
```

### 7. Implement the Send Method

```go
func (s *MyService) Send(message string) error {
    // Use s.options.apiKey, s.options.endpoint, etc.
    // Implementation here
    return nil
}
```

### 8. Usage

```go
// Build options
opts := integrations.NewMyServiceOptions().
    APIKey("your-api-key").
    Endpoint("https://api.example.com").
    Build()

// Create client with options
service, err := integrations.NewMyService(opts)
if err != nil {
    log.Fatal(err)
}

// Use the client
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

Validation is automatically performed when calling `New<Integration>(opts)`, ensuring invalid configurations are caught before the client is created.

## Configuration

### Using the Options Builder Pattern (Recommended)

```go
opts := integrations.NewSMTPOptions().
    Server("smtp.gmail.com").
    Port(587).
    Username("user@example.com").
    Password("password").
    From("from@example.com").
    To("to@example.com").
    Build()

smtp, err := integrations.NewSMTP(opts)
```
### Environment Variables

You can load configuration from environment variables:

```go
import "os"

opts := integrations.NewSMTPOptions().
    Server(os.Getenv("SMTP_SERVER")).
    Port(587).
    Username(os.Getenv("SMTP_USERNAME")).
    Password(os.Getenv("SMTP_PASSWORD")).
    From(os.Getenv("SMTP_FROM")).
    To(os.Getenv("SMTP_TO")).
    Build()

smtp, err := integrations.NewSMTP(opts)
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

The options builder pattern provides clear error handling:

```go
// Build options (no error here)
opts := integrations.NewSMTPOptions().
    Server("smtp.gmail.com").
    Port(587).
    // Missing required fields...
    Build()

// Validation happens when creating the client
smtp, err := integrations.NewSMTP(opts)
if err != nil {
    // Validation error caught at client creation time
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

Contributions are welcome! When adding new integrations, please follow the options builder pattern demonstrated in this repository.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Follow the options builder pattern (see "Creating a New Integration" above)
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
feat(smtp): add options builder pattern implementation
feat(telegram): add options builder with method chaining
fix(slack): correct webhook payload formatting
docs(readme): update usage examples with options builder pattern
refactor(pushover): migrate to options builder pattern
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

## Benefits of the Options Builder Pattern

### Clean Syntax
Build options separately from client creation:
```go
opts := integrations.NewSMTPOptions().
    Server("smtp.gmail.com").
    Port(587).
    Build()

smtp, err := integrations.NewSMTP(opts)
```

### Separation of Concerns
Options building is completely separate from client creation, following the same pattern as MongoDB's official driver.

### Type Safety
Compile-time type checking prevents configuration errors.

### Flexibility
Configure only the options you need. Method chaining is optional.

### Validation
Built-in validation when creating the client ensures configurations are correct before use, catching errors early.

### Extensibility
Adding new builder methods doesn't break existing code. Simply add new chainable methods to the options builder.

### Readability
Self-documenting fluent API makes code easy to read and understand:
```go
// Clear and readable - MongoDB style
opts := integrations.NewSMTPOptions().
    Server("smtp.gmail.com").
    Port(587).
    Username("user@example.com").
    Password("password").
    From("from@example.com").
    To("to@example.com").
    Build()

smtp, err := integrations.NewSMTP(opts)
```

## Support

- Issues: [GitHub Issues](https://github.com/uug-ai/integrations/issues)
- Discussions: [GitHub Discussions](https://github.com/uug-ai/integrations/discussions)
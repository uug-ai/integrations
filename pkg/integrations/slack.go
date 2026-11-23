package integrations

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/uug-ai/integrations/pkg/message"
	"github.com/uug-ai/models/pkg/models"
)

const slackMaxVideoBytes int64 = 1_000_000_000

// slackSetupError is a Slack delivery that cannot succeed with this
// configuration or recording, so retrying it does not help.
type slackSetupError string

func (err slackSetupError) Error() string          { return string(err) }
func (err slackSetupError) DeliveryRejected() bool { return true }

// SlackOptions configures Slack.
type SlackOptions struct {
	// DeliveryMode is the delivery mode: "webhook" or "bot"; empty uses the fields that are set.
	DeliveryMode string `json:"deliveryMode,omitempty" validate:"omitempty,oneof=webhook bot"`
	// Token is the bot token.
	Token string `json:"token,omitempty"`
	// Channel is the channel ID the bot posts to.
	Channel string `json:"channel,omitempty"`
	// Hook is the incoming webhook URL.
	Hook string `json:"hook,omitempty" validate:"omitempty,url"`
	// Username is the name webhook messages are posted as.
	Username string `json:"username,omitempty"`
}

// SlackOptionsBuilder builds SlackOptions.
type SlackOptionsBuilder struct {
	options *SlackOptions
}

// NewSlackOptions starts building SlackOptions.
func NewSlackOptions() *SlackOptionsBuilder {
	return &SlackOptionsBuilder{options: &SlackOptions{}}
}

// SetDeliveryMode sets the delivery mode: "webhook" or "bot"; empty uses the fields that are set.
func (b *SlackOptionsBuilder) SetDeliveryMode(deliveryMode string) *SlackOptionsBuilder {
	b.options.DeliveryMode = deliveryMode
	return b
}

// SetToken sets the bot token.
func (b *SlackOptionsBuilder) SetToken(token string) *SlackOptionsBuilder {
	b.options.Token = token
	return b
}

// SetChannel sets the channel ID the bot posts to.
func (b *SlackOptionsBuilder) SetChannel(channel string) *SlackOptionsBuilder {
	b.options.Channel = channel
	return b
}

// SetHook sets the incoming webhook URL.
func (b *SlackOptionsBuilder) SetHook(hook string) *SlackOptionsBuilder {
	b.options.Hook = hook
	return b
}

// SetUsername sets the name webhook messages are posted as.
func (b *SlackOptionsBuilder) SetUsername(username string) *SlackOptionsBuilder {
	b.options.Username = username
	return b
}

// Build returns the configured SlackOptions.
func (b *SlackOptionsBuilder) Build() *SlackOptions {
	return b.options
}

// Slack sends messages to a Slack channel, through a webhook or a bot.
type Slack struct {
	options *SlackOptions
}

// NewSlack validates opts and returns a Slack sender that uses a copy of them.
func NewSlack(opts *SlackOptions) (*Slack, error) {
	if err := validateOptions("slack", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Slack{options: &options}, nil
}

func (s *Slack) Send(m message.Message) error {
	return s.SendWithContext(context.Background(), m)
}

func (s *Slack) SendWithContext(ctx context.Context, m message.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}
	api := slack.New(s.options.Token, slack.OptionHTTPClient(client))
	return s.send(ctx, m, client, api)
}

func (s *Slack) send(ctx context.Context, m message.Message, client *http.Client, api *slack.Client) error {
	text := m.Body
	attachments := []slack.Attachment{}
	if len(m.Media) > 0 {
		if m.Media[0].Url != "" {
			text += "\r\n" + m.Media[0].Url
		}
		if m.Media[0].ThumbnailUrl != "" {
			attachments = append(attachments, slack.Attachment{
				Color:    "good",
				ImageURL: m.Media[0].ThumbnailUrl,
			})
		}
	}

	mode, err := (models.Slack{
		DeliveryMode: s.options.DeliveryMode, Apikey: s.options.Token, Channel: s.options.Channel,
	}).ResolveDeliveryMode()
	if err != nil {
		return slackSetupError(err.Error())
	}
	if mode == models.SlackDeliveryModeBot {
		if s.options.Token == "" || s.options.Channel == "" {
			return slackSetupError("Slack bot delivery requires both a bot token and channel ID")
		}
		if len(m.Media) > 0 && m.Media[0].Type == "video" {
			return s.uploadVideo(ctx, m, client, api)
		}
		_, _, err := api.PostMessageContext(ctx, s.options.Channel,
			slack.MsgOptionText(text, false), slack.MsgOptionAttachments(attachments...))
		return err
	}

	if s.options.Hook == "" {
		return slackSetupError("Slack webhook delivery requires a webhook URL")
	}
	msg := slack.WebhookMessage{
		Username:    s.options.Username,
		Text:        text,
		Attachments: attachments,
	}
	return slack.PostWebhookCustomHTTPContext(ctx, s.options.Hook, client, &msg)
}

func (s *Slack) uploadVideo(ctx context.Context, m message.Message, client *http.Client, api *slack.Client) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, m.Media[0].Url, nil)
	if err != nil {
		return slackSetupError("invalid Slack video download URL")
	}
	response, err := client.Do(request)
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return fmt.Errorf("download Slack video: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download Slack video: %w", slack.StatusCodeError{
			Code: response.StatusCode, Status: http.StatusText(response.StatusCode),
		})
	}
	if response.ContentLength > slackMaxVideoBytes {
		return slackSetupError("Slack video exceeds the 1 GB file limit")
	}

	// Slack requires the exact byte length before issuing an upload URL.
	file, err := os.CreateTemp("", "hub-slack-video-*.mp4")
	if err != nil {
		return fmt.Errorf("create Slack video temporary file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Slack video temporary file close failed: %v", err)
		}
		if err := os.Remove(file.Name()); err != nil {
			log.Printf("Slack video temporary file cleanup failed: %v", err)
		}
	}()

	size, err := copySlackVideo(file, response.Body)
	if err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind Slack video: %w", err)
	}
	title := m.Title
	if title == "" {
		title = "Recording"
	}
	_, err = api.UploadFileContext(ctx, slack.UploadFileParameters{
		Reader:         file,
		FileSize:       int(size),
		Filename:       "recording.mp4",
		Title:          title,
		InitialComment: m.Body,
		Channel:        s.options.Channel,
	})
	if err != nil {
		return fmt.Errorf("upload Slack video: %w", err)
	}
	return nil
}

func copySlackVideo(dst io.Writer, src io.Reader) (int64, error) {
	size, err := io.Copy(dst, io.LimitReader(src, slackMaxVideoBytes+1))
	if err != nil {
		return 0, fmt.Errorf("read Slack video: %w", err)
	}
	if size == 0 {
		return 0, slackSetupError("Slack video is empty")
	}
	if size > slackMaxVideoBytes {
		return 0, slackSetupError("Slack video exceeds the 1 GB file limit")
	}
	return size, nil
}

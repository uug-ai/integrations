package integrations

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/uug-ai/integrations/pkg/message"
)

// TelegramOptions configures Telegram.
type TelegramOptions struct {
	// Token is the bot token.
	Token string `json:"token,omitempty" validate:"required"`
	// Channel is the channel name, or c<id>_<hash> for a private channel.
	Channel string `json:"channel,omitempty" validate:"required"`
}

// TelegramOptionsBuilder builds TelegramOptions.
type TelegramOptionsBuilder struct {
	options *TelegramOptions
}

// NewTelegramOptions starts building TelegramOptions.
func NewTelegramOptions() *TelegramOptionsBuilder {
	return &TelegramOptionsBuilder{options: &TelegramOptions{}}
}

// SetToken sets the bot token.
func (b *TelegramOptionsBuilder) SetToken(token string) *TelegramOptionsBuilder {
	b.options.Token = token
	return b
}

// SetChannel sets the channel name, or c<id>_<hash> for a private channel.
func (b *TelegramOptionsBuilder) SetChannel(channel string) *TelegramOptionsBuilder {
	b.options.Channel = channel
	return b
}

// Build returns the configured TelegramOptions.
func (b *TelegramOptionsBuilder) Build() *TelegramOptions {
	return b.options
}

// Telegram posts the message to a Telegram channel with a bot.
type Telegram struct {
	options *TelegramOptions
}

// NewTelegram validates opts and returns a Telegram sender that uses a copy of them.
func NewTelegram(opts *TelegramOptions) (*Telegram, error) {
	if err := validateOptions("telegram", opts); err != nil {
		return nil, err
	}
	options := *opts
	return &Telegram{options: &options}, nil
}

func (t *Telegram) Send(message message.Message) bool {

	channelName := t.options.Channel //"c1375189391_8694429167782276799"
	token := t.options.Token         //"592498002:AAHYGK-EEUXV3oFtf3mVUJEPWxCYfNkXdC0"

	// Create a bot with BotFather
	// /newbot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return false
	}
	bot.Debug = false

	// Get channel id through web interface.
	// e.g. https://web.telegram.org/#/im?p=c1375189391_8694429167782276799
	// Provide (private channel): c1375189391_8694429167782276799

	isPrivate := (channelName[:1] == "c")
	if isPrivate {
		channelId := strings.Split(channelName, "_")[0]
		channelName = strings.Replace(channelId, "c", "-100", 1)
	}

	// Shorten url
	url := ""
	if len(message.Media) > 0 {
		longUrl := message.Media[0].Url
		url = longUrl
		//provider := "tinyurl"
		//shortenedUrl, err := shorturl.Shorten(longUrl, provider)
		//if err == nil {
		// url = string(shortenedUrl)
		//}
	}

	text := message.Body
	if url != "" {
		text = text + "\r\n" + url
	}

	msg := tgbotapi.NewMessageToChannel(channelName, text)
	bot.Send(msg)

	return true
}

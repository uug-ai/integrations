package pkg

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Telegram struct {
	Token   string `json:"token,omitempty"`
	Channel string `json:"channel,omitempty"`
}

func (t Telegram) Send(message Message) bool {

	channelName := t.Channel //"c1375189391_8694429167782276799"
	token := t.Token         //"592498002:AAHYGK-EEUXV3oFtf3mVUJEPWxCYfNkXdC0"

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

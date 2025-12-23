package channels

import (
	tw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	message "github.com/uug-ai/hub-pipeline-notification/message"
)

type Twitter struct{}

func (twitter Twitter) SendNotification(message message.Message) bool {

	// Twitter client
	config := oauth1.NewConfig("consumerKey", "pTMuNmxrTlG7IlmqXEzeqjfx7nUFPiwUXsXNUhNOWlYkHxLjXz")
	token := oauth1.NewToken("accessToken", "	B1HlS7dtnRtoEl7XJs6NrzOSazlehgHUosYtLNJliYFzW")
	httpClient := config.Client(oauth1.NoContext, token)
	client := tw.NewClient(httpClient)

	// Send a Tweet.
	client.Statuses.Update("just setting up my twttr", nil)

	return true
}

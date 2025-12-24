package integrations

import (
	tw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/uug-ai/models/pkg/models"
)

type Twitter struct{}

func (twitter Twitter) SendNotification(message models.Message) bool {

	// Twitter client
	config := oauth1.NewConfig("consumerKey", "pTMuNmxrTlG7IlmqXEzeqjfx7nUFPiwUXsXNUhNOWlYkHxLjXz")
	token := oauth1.NewToken("accessToken", "	B1HlS7dtnRtoEl7XJs6NrzOSazlehgHUosYtLNJliYFzW")
	httpClient := config.Client(oauth1.NoContext, token)
	client := tw.NewClient(httpClient)

	// Send a Tweet.
	client.Statuses.Update("just setting up my twttr", nil)

	return true
}

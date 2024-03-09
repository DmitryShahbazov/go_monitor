package scrapers

import (
	"github.com/nicklaw5/helix"
	"youtube_monitor/settings"
)

func getClientApi() *helix.Client {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     settings.TwitchClientID,
		ClientSecret: settings.TwitchClientSecret,
	})
	if err != nil {
		// handle error
	}

	resp, err := client.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		// handle error
	}

	client.SetAppAccessToken(resp.Data.AccessToken)
	return client
}

func CheckIsLiveTwitch(channel, nickname string, results chan<- LiveCheckResult) {
	client := getClientApi()

	resp, err := client.GetStreams(&helix.StreamsParams{
		First:      10,
		UserLogins: []string{channel},
	})
	if err != nil {
		results <- LiveCheckResult{0, "", err, ""}
	}

	if len(resp.Data.Streams) > 0 {
		URL := "https://twitch.tv/" + channel
		results <- LiveCheckResult{int32(resp.Data.Streams[0].ViewerCount), nickname, nil, URL}
	}

}

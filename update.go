package main

import (
	aw "github.com/deanishe/awgo"
	"github.com/slack-go/slack"
)

func updateChannels() {
	wf.NewItem("Update Channels").Valid(true)

	c := aw.NewCache(cache_dir)
	cfg := aw.NewConfig()
	token := cfg.Get("SLACK_TOKEN")
	api := slack.New(token)

	channels, err := getConversations(api, []slack.Channel{}, "")
	team, err_team := api.GetTeamInfo()
	if err != nil || err_team != nil {
		wf.Warn("Error", "Error occurred in Slack API ")
	}

	all_channels := make([]Channel, 0)
	for _, channel := range channels {
		all_channels = append(all_channels, Channel{
			Name:   channel.Name,
			ID:     channel.ID,
			TeamID: team.ID,
		})
	}

	c.StoreJSON(cache_file, all_channels)
	wf.SendFeedback()
}

func getConversations(api *slack.Client, channels []slack.Channel, cursor string) ([]slack.Channel, error) {
	params := slack.GetConversationsParameters{
		Limit: 200,
		ExcludeArchived: "true",
		Cursor: cursor,
	}
	next_channels, next_cursor, err_channels := api.GetConversations(&params)
	if err_channels != nil {
		return nil, err_channels
	}

	channels = append(channels, next_channels...)

	if next_cursor == "" {
		return channels, nil
	}
	
	return getConversations(api, channels, next_cursor)
}

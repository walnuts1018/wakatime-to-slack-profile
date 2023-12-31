package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
)

type client struct {
	slackClient *slack.Client
}

func NewClient() *client {
	return &client{
		slackClient: slack.New(config.Config.SlackAccessToken),
	}
}

func (c *client) SetUserCustomStatus(emoji string) error {
	if !(strings.HasPrefix(emoji, ":") && strings.HasSuffix(emoji, ":")) {
		emoji = ":" + emoji + ":"
	}
	err := c.slackClient.SetUserCustomStatus("", emoji, 0)
	if err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	return nil
}

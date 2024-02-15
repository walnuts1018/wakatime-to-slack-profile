package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
)

type client struct {
	slackClient *slack.Client
}

func NewClient(cfg config.Config) domain.SlackClient {
	return &client{
		slackClient: slack.New(cfg.SlackAccessToken),
	}
}

func (c *client) SetUserCustomStatus(emoji string, text string) error {
	if !(strings.HasPrefix(emoji, ":") && strings.HasSuffix(emoji, ":")) {
		emoji = ":" + emoji + ":"
	}
	if text != "" {
		text = fmt.Sprintf("now coding: %v", text)
	}
	err := c.slackClient.SetUserCustomStatus(text, emoji, 0)
	if err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	return nil
}

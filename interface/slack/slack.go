package slack

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

const slackRedirectPath = "/slack_callback"

type Client struct {
	cfg config.Config
}

func NewClient(cfg config.Config) (Client, error) {
	return Client{cfg: cfg}, nil
}

func (c Client) GetUser(id string, token string) (model.SlackUser, error) {
	user, err := slack.New(token).GetUserInfo(id)
	if err != nil {
		return model.SlackUser{}, fmt.Errorf("error getting user: %w", err)
	}

	return model.SlackUser{
		ID:     user.ID,
		Name:   user.Name,
		TeamID: user.TeamID,
	}, nil
}

func (c Client) SetCustomStatus(id, token, emoji, text string) error {
	err := slack.New(token).SetUserCustomStatus(text, emoji, 0)
	if err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	return nil
}

func (c Client) GetEmojis(teamID, token string) (model.SlackEmojis, error) {

	es, err := slack.New(token).GetEmoji()
	if err != nil {
		return model.SlackEmojis{}, fmt.Errorf("error getting emojis: %w", err)
	}

	var emojis map[string]struct{}

	for k := range es {
		emojis[k] = struct{}{}
	}

	return model.SlackEmojis{
		TeamID: teamID,
		Emojis: emojis,
	}, nil
}

func (c Client) GetOAuthV2Response(code string) (string, string, error) {
	redirectURL, err := url.JoinPath(c.cfg.ServerURL, slackRedirectPath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to join url: %v", err), "server_url", c.cfg.ServerURL, "path", slackRedirectPath)
		os.Exit(1)
	}

	resp, err := slack.GetOAuthV2Response(
		http.DefaultClient,
		c.cfg.SlackClientID,
		c.cfg.SlackClientSecret,
		code,
		redirectURL,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to get oauth response: %w", err)
	}

	return resp.AccessToken, resp.AuthedUser.ID, nil
}

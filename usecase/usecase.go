package usecase

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
)

type Usecase struct {
	wakatimeClient domain.WakatimeClient
	tokenStore     domain.TokenStore
	slackClient    domain.SlackClient
	emojiOverides  map[string]string
	lastLanguage   *string
}

func NewUsecase(wakatimeClient domain.WakatimeClient, tokenStore domain.TokenStore, slackClient domain.SlackClient) *Usecase {
	var emojis map[string]string
	b, err := os.ReadFile("emoji.json")
	if err != nil {
		slog.Warn("emoji.json is not found")
	} else {
		json.Unmarshal(b, &emojis)
	}
	return &Usecase{
		wakatimeClient: wakatimeClient,
		tokenStore:     tokenStore,
		slackClient:    slackClient,
		emojiOverides:  emojis,
	}
}

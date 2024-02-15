package usecase

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
)

type Usecase struct {
	wakatimeClient domain.WakatimeClient
	db             domain.DB
	slackClient    domain.SlackClient
	emojiOverides  map[string]string
	lastLanguage   *string
}

func NewUsecase(wakatimeClient domain.WakatimeClient, db domain.DB, slackClient domain.SlackClient, emojiOverides map[string]string) *Usecase {
	emojis := map[string]string{}

	b, err := os.ReadFile("emoji.json")
	if err != nil {
		slog.Warn("emoji.json is not found")
	} else {
		json.Unmarshal(b, &emojis)
	}

	for k, v := range emojiOverides {
		emojis[k] = v
	}

	return &Usecase{
		wakatimeClient: wakatimeClient,
		db:             db,
		slackClient:    slackClient,
		emojiOverides:  emojis,
	}
}

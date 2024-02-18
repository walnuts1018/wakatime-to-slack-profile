package repository

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

type SlackEmojis interface {
	GetEmojis(teamID string) (model.SlackEmojis, error)
}

type SlackEmojisCache interface {
	CacheEmojis(emojis model.SlackEmojis) error
	IsEmojiExist(teamID, emoji string) (bool, error)
	IsCached(teamID string) (bool, synchro.Time[tz.AsiaTokyo], error)
}

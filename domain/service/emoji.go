package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/repository"
)

type SlackEmojiService struct {
	SlackEmojisRepository repository.SlackEmojis
	SlackEmojisCache      repository.SlackEmojisCache

	emojiCacheDuration time.Duration
}

func NewSlackEmojiService(cfg config.Config, slackEmojisRepository repository.SlackEmojis, slackEmojisCache repository.SlackEmojisCache) *SlackEmojiService {
	return &SlackEmojiService{
		emojiCacheDuration:    cfg.EmojiCacheDuration,
		SlackEmojisRepository: slackEmojisRepository,
		SlackEmojisCache:      slackEmojisCache,
	}
}

func (s *SlackEmojiService) IsEmojiExist(teamID, emoji string) (bool, error) {
	cached, lastupdate, err := s.SlackEmojisCache.IsCached(teamID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check if emojis are cached: %v", err), slog.String("teamID", teamID))
		cached = false
	}

	now := synchro.Now[tz.AsiaTokyo]()

	if cached && now.Sub(lastupdate) < s.emojiCacheDuration {
		return s.SlackEmojisCache.IsEmojiExist(teamID, emoji)
	} else {
		emojis, err := s.SlackEmojisRepository.GetEmojis(teamID)
		if err != nil {
			return false, fmt.Errorf("failed to get emojis from slack: %w", err)
		}

		if err := s.SlackEmojisCache.CacheEmojis(emojis); err != nil {
			slog.Error(fmt.Sprintf("failed to cache emojis: %v", err), slog.String("teamID", teamID))
		}

		_, ok := emojis.Emojis[emoji]
		return ok, nil
	}
}

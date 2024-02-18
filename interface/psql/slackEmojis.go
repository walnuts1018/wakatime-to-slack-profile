package psql

import (
	"fmt"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
	"gorm.io/gorm"
)

func (c Client) CacheEmojis(emojis model.SlackEmojis) error {
	var obj []SlackEmoji
	for emoji := range emojis.Emojis {
		obj = append(obj, SlackEmoji{
			TeamID:      emojis.TeamID,
			Emoji:       emoji,
			LastUpdated: synchro.Now[tz.AsiaTokyo](),
		})
	}

	result := c.db.Create(&obj)
	if result.Error != nil {
		return fmt.Errorf("failed to cache emojis: %w", result.Error)
	}

	return nil

}

// IsEmojiExist checks if the emoji exists in the cache, returns if it exists, the last updated time and an error
func (c Client) IsEmojiExist(teamID, emoji string) (bool, synchro.Time[tz.AsiaTokyo], error) {
	var slackEmoji SlackEmoji
	result := c.db.Where("team_id = ? AND emoji = ?", teamID, emoji).First(&slackEmoji)
	if result.Error == gorm.ErrRecordNotFound {
		return false, synchro.Time[tz.AsiaTokyo]{}, nil
	}

	if result.Error != nil {
		return false, synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("failed to check if emoji exists: %w", result.Error)
	}
	return true, slackEmoji.LastUpdated, nil
}

func (c Client) IsCached(teamID string) (bool, error) {
	var emoji SlackEmoji
	result := c.db.Where("team_id = ?", teamID).First(&emoji)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	}

	if result.Error != nil {
		return false, fmt.Errorf("failed to check if emojis are cached: %w", result.Error)
	}
	return true, nil
}

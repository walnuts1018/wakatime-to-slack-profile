package psql

import (
	"fmt"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

func (c Client) UpdateSlackToken(slackToken model.SlackToken) error {
	result := c.db.Create(&slackToken)
	if result.Error != nil {
		return fmt.Errorf("failed to update slack token: %w", result.Error)
	}
	return nil
}

func (c Client) GetSlackToken(slackUserID string) (SlackToken, error) {
	var token SlackToken
	result := c.db.Where("slack_user_id = ?", slackUserID).First(&token)
	if result.Error != nil {
		return SlackToken{}, fmt.Errorf("failed to get slack token: %w", result.Error)
	}
	return token, nil
}

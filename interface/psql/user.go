package psql

import "fmt"

func (c Client) AddUser(user User) error {
	result := c.db.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to add user: %w", result.Error)
	}
	return nil
}

func (c Client) GetUser(userID string) (User, error) {
	var user User
	result := c.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return User{}, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return user, nil
}

func (c Client) UpdateUser(user User) error {
	result := c.db.Save(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

func (c Client) AddOrUpdateEmoji(userID string, emojis map[string]string) error {
	for language, emoji := range emojis {
		c.db.Create(&CustomEmojis{
			UserID:   userID,
			Language: language,
			Emoji:    emoji,
		})
	}

	return nil
}

func (c Client) RemoveEmoji(userID, language string) error {
	result := c.db.Where("user_id = ? AND language = ?", userID, language).Delete(&CustomEmojis{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete emoji: %w", result.Error)
	}
	return nil
}
func (c Client) GetEmoji(userID, language string) (string, error) {
	var emoji CustomEmojis
	result := c.db.Where("user_id = ? AND language = ?", userID, language).First(&emoji)
	if result.Error != nil {
		return "", fmt.Errorf("failed to get emoji: %w", result.Error)
	}
	return emoji.Emoji, nil
}

func (c Client) GetAllEmojis(userID string) (map[string]string, error) {
	var emojis []CustomEmojis
	result := c.db.Where("user_id = ?", userID).Find(&emojis)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all emojis: %w", result.Error)
	}

	emojiMap := make(map[string]string)
	for _, emoji := range emojis {
		emojiMap[emoji.Language] = emoji.Emoji
	}
	return emojiMap, nil
}

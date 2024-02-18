package model

type User struct {
	ID             string            `json:"id"`
	Username       string            `json:"username"`
	SlackUserID    string            `json:"slack_user_id"`
	WakatimeUserID string            `json:"wakatime_user_id"`
	emojis         map[string]string // language: emojiID
}

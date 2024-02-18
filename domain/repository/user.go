package repository

import "github.com/walnuts1018/wakatime-to-slack-profile/domain/model"

type User interface {
	AddUser(user model.User) error
	GetUser(userID string) (model.User, error)
	UpdateUser(user model.User) error

	AddOrUpdateEmoji(userID string, emojis map[string]string) error
	RemoveEmoji(userID string, language string) error
	GetEmoji(userID, language string) (string, error)
	GetAllEmojis(userID string) (map[string]string, error)
}

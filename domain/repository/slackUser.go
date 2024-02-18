package repository

import "github.com/walnuts1018/wakatime-to-slack-profile/domain/model"

type SlackUser interface {
	GetUser(id string, token string) (model.SlackUser, error)
	SetCustomStatus(id, token, emoji, text string) error
}

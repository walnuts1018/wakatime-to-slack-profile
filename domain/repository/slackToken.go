package repository

import (
	"errors"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

type SlackToken interface {
	UpdateSlackToken(slackToken model.SlackToken) error
	GetSlackToken(slackUserID string) (model.SlackToken, error)
}

var ErrorSlackTokenNotFound = errors.New("slack token not found")

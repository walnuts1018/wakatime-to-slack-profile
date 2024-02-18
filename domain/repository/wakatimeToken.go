package repository

import (
	"errors"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

type WakatimeToken interface {
	UpdateWakatimeToken(token model.WakatimeToken) error
	GetWakatimeToken(wakatimeUserID string) (model.WakatimeToken, error)
}

var ErrorWakatimeTokenNotFound = errors.New("wakatime token not found")

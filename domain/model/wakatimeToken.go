package model

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type WakatimeToken struct {
	WakatimeUserID string                     `json:"user_id"`
	AccessToken    string                     `json:"access_token"`
	RefreshToken   string                     `json:"refresh_token"`
	ExpiresAt      synchro.Time[tz.AsiaTokyo] `json:"expires_at"`
	UpdatedAt      synchro.Time[tz.AsiaTokyo] `json:"updated_at"`
}

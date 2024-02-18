package model

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

// 今のところAccessTokenが無期限で使えるので、refresh_tokenは未実装
type SlackToken struct {
	SlackUserID  string                     `json:"user_id"`
	AccessToken  string                     `json:"access_token"`
	RefreshToken string                     `json:"refresh_token"` // 未実装
	ExpiresAt    synchro.Time[tz.AsiaTokyo] `json:"expires_at"`    // 未実装
	UpdatedAt    synchro.Time[tz.AsiaTokyo] `json:"updated_at"`    // 未実装
}

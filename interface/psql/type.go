package psql

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type CustomEmojis struct {
	UserID   string `gorm:"primaryKey"`
	Language string `gorm:"primaryKey"`
	Emoji    string
}

type SlackToken struct {
	SlackUserID  string `gorm:"primaryKey"`
	AccessToken  string
	RefreshToken string
	ExpiresAt    synchro.Time[tz.AsiaTokyo]
	UpdatedAt    synchro.Time[tz.AsiaTokyo]
}

type WakatimeToken struct {
	WakatimeUserID string `gorm:"primaryKey"`
	AccessToken    string
	RefreshToken   string
	ExpiresAt      synchro.Time[tz.AsiaTokyo]
	UpdatedAt      synchro.Time[tz.AsiaTokyo]
}

type User struct {
	ID             string `gorm:"primaryKey"`
	Username       string
	SlackUserID    string
	WakatimeUserID string
}

type SlackEmoji struct {
	TeamID      string `gorm:"primaryKey"`
	Emoji       string `gorm:"primaryKey"`
	LastUpdated synchro.Time[tz.AsiaTokyo]
}

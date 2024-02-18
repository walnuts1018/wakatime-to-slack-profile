package model

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type SlackEmojis struct {
	TeamID      string
	Emojis      map[string]struct{}
	LastUpdated synchro.Time[tz.AsiaTokyo]
}

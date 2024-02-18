package model

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type WakatimeUser struct {
	ID   string `json:"id"`
	Name string `json:"username"`
}

type WakatimeLastActivity struct {
	Language string
	Project  string
	Start    synchro.Time[tz.AsiaTokyo]
	End      synchro.Time[tz.AsiaTokyo]
}

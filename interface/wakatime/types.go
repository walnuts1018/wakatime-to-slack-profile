package wakatime

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type durationsResponce struct {
	Data []struct {
		Duration float64 `json:"duration"`
		Language string  `json:"language"`
		Project  string  `json:"project"`
		Time     float64 `json:"time"`
	} `json:"data"`
	Start    synchro.Time[tz.AsiaTokyo] `json:"start"`
	End      synchro.Time[tz.AsiaTokyo] `json:"end"`
	Timezone string                     `json:"timezone"`
}

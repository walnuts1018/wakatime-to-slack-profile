package repository

import (
	"net/http"

	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

type WakatimeUser interface {
	GetUser(wakatimeUserID string, oauth2Client *http.Client) (model.WakatimeUser, error)
	GetCurrentUser(oauth2Client *http.Client) (model.WakatimeUser, error)
	GetLastActivity(wakatimeUserID string, oauth2Client *http.Client) (model.WakatimeLastActivity, error)
}

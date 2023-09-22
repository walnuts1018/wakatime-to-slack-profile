package domain

import (
	"context"
	"time"
)

type OAuth2Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WakatimeClient interface {
	Auth(state string) string
	Callback(ctx context.Context, code string) (OAuth2Token, error)
	SetToken(ctx context.Context, tokenStore TokenStore) error
	ListLanguages(ctx context.Context) ([]Language, error)
	NowLanguage(ctx context.Context) (string, error)
}

type Language struct {
	Id         string `json:"id"`          //unique id of this language
	Name       string `json:"name"`        //human readable name of this language
	Color      string `json:"color"`       //hex color code, used when displaying this language on WakaTime charts
	IsVerified bool   `json:"is_verified"` //whether this language is verified, by GitHub’s linguist or manually by WakaTime admins
	CreatedAt  string `json:"created_at"`  //time when this language was created in ISO 8601 format
	ModifiedAt string `json:"modified_at"` //time when this language was last modified in ISO 8601 format
}

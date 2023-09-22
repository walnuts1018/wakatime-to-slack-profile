package usecase

import "github.com/walnuts1018/wakatime-to-slack-profile/domain"

type Usecase struct {
	wakatimeClient domain.WakatimeClient
	tokenStore     domain.TokenStore
}

func NewUsecase(wakatimeClient domain.WakatimeClient, tokenStore domain.TokenStore) *Usecase {
	return &Usecase{
		wakatimeClient: wakatimeClient,
		tokenStore:     tokenStore,
	}
}

package wakatime

import (
	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/timeJST"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	src        oauth2.TokenSource
	tokenStore domain.TokenStore
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	token := domain.OAuth2Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		UpdatedAt:    timeJST.Now(),
	}

	err = s.tokenStore.UpdateOAuth2Token(token)
	if err != nil {
		return t, err
	}
	return t, nil
}

package usecase

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	src        oauth2.TokenSource
	tokenStore domain.DB
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	token := domain.OAuth2Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       synchro.In[tz.AsiaTokyo](t.Expiry),
		UpdatedAt:    synchro.Now[tz.AsiaTokyo](),
	}

	err = s.tokenStore.UpdateOAuth2Token(token)
	if err != nil {
		return t, err
	}
	return t, nil
}

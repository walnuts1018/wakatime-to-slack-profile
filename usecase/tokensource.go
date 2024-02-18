package usecase

import (
	"fmt"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/repository"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	src               oauth2.TokenSource
	wakatimeTokenRepo repository.WakatimeToken
	wakatimeUserID    string
}

func (s *tokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	token := model.WakatimeToken{
		WakatimeUserID: s.wakatimeUserID,
		AccessToken:    t.AccessToken,
		RefreshToken:   t.RefreshToken,
		ExpiresAt:      synchro.In[tz.AsiaTokyo](t.Expiry),
		UpdatedAt:      synchro.Now[tz.AsiaTokyo](),
	}

	if err := s.wakatimeTokenRepo.UpdateWakatimeToken(token); err != nil {
		return t, fmt.Errorf("failed to update wakatime token: %w", err)
	}
	return t, nil
}

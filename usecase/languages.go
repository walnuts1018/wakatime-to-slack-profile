package usecase

import (
	"context"
	"fmt"
)

func (u *Usecase) SetLanguage(ctx context.Context) error {
	language, err := u.wakatimeClient.NowLanguage(ctx)
	if err != nil {
		return fmt.Errorf("failed to get now language: %w", err)
	}

	if u.lastLanguage != nil {
		if *u.lastLanguage == language {
			return nil
		}
	}

	err = u.SetUserCustomStatus(language)
	if err != nil {
		return fmt.Errorf("failed to set user custom status: %w", err)
	}
	u.lastLanguage = &language
	return nil
}

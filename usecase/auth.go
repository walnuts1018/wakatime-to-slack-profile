package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
)

func (u *Usecase) SignIn() (string, string, error) {
	state, err := randStr(64)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random string: %w", err)
	}
	redirect := u.wakatimeClient.Auth(state)
	return state, redirect, nil
}

func (u *Usecase) Callback(ctx context.Context, code string) error {
	token, err := u.wakatimeClient.Callback(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get oauth2 config: %w", err)
	}

	err = u.tokenStore.SaveOAuth2Token(token)
	if err != nil {
		return fmt.Errorf("failed to save oauth2 token: %w", err)
	}

	return nil
}

func (u *Usecase) SetToken(ctx context.Context) error {
	return u.wakatimeClient.SetToken(ctx, u.tokenStore)
}

func randStr(n int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

package wakatime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/timeJST"
	"golang.org/x/oauth2"
)

const (
	AuthEndpoint  = "https://wakatime.com/oauth/authorize"
	TokenEndpoint = "https://wakatime.com/oauth/token"
)

var (
	scopes = []string{
		"read_stats",
	}
)

type client struct {
	cfg     *oauth2.Config
	wclient *http.Client
}

func NewOauth2Client() domain.WakatimeClient {
	return &client{
		cfg: &oauth2.Config{
			ClientID:     config.Config.WakatimeAppID,
			ClientSecret: config.Config.WakatimeAppSecret,
			Endpoint:     oauth2.Endpoint{AuthURL: AuthEndpoint, TokenURL: TokenEndpoint},
			Scopes:       scopes,
		},
	}
}

func (c *client) Auth(state string) string {
	url := c.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

func (c *client) Callback(ctx context.Context, code string) (domain.OAuth2Token, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return domain.OAuth2Token{}, err
	}

	cfg := domain.OAuth2Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		CreatedAt:    timeJST.Now(),
		UpdatedAt:    timeJST.Now(),
	}
	return cfg, nil
}

func (c *client) SetToken(ctx context.Context, tokenStore domain.TokenStore) error {
	token, err := tokenStore.GetOAuth2Token()
	if err != nil {
		return fmt.Errorf("failed to get oauth2 token: %w", err)
	}
	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    "bearer",
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}

	oldTokenSource := c.cfg.TokenSource(ctx, oauthToken)
	mySrc := &tokenSource{
		src:        oldTokenSource,
		tokenStore: tokenStore,
	}

	reuseSrc := oauth2.ReuseTokenSource(oauthToken, mySrc)
	c.wclient = oauth2.NewClient(ctx, reuseSrc)
	return nil
}

func (c *client) Languages(ctx context.Context) ([]domain.Language, error) {
	if c.wclient == nil {
		return nil, fmt.Errorf("client is not set")
	}
	resp, err := c.wclient.Get("https://wakatime.com/api/v1/program_languages")
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get languages: %v", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var languages []domain.Language
	err = json.Unmarshal(raw, &languages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return languages, nil

}

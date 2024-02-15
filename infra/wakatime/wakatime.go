package wakatime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain"
	"golang.org/x/oauth2"
)

const (
	AuthEndpoint  = "https://wakatime.com/oauth/authorize"
	TokenEndpoint = "https://wakatime.com/oauth/token"
)

var (
	scopes = []string{
		"read_stats",
		"read_logged_time",
	}
)

type client struct {
	cfg     *oauth2.Config
	wclient *http.Client
}

func NewOauth2Client(cfg config.Config) (domain.WakatimeClient, error) {
	url, err := url.JoinPath(cfg.ServerURL, "callback")
	if err != nil {
		return nil, fmt.Errorf("failed to join url: %w", err)
	}

	return &client{
		cfg: &oauth2.Config{
			ClientID:     cfg.WakatimeAppID,
			ClientSecret: cfg.WakatimeAppSecret,
			Endpoint:     oauth2.Endpoint{AuthURL: AuthEndpoint, TokenURL: TokenEndpoint},
			RedirectURL:  url,
			Scopes:       scopes,
		},
	}, nil
}

func (c *client) Auth(state string) string {
	url := c.cfg.AuthCodeURL(state)
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
		Expiry:       synchro.In[tz.AsiaTokyo](token.Expiry),
		CreatedAt:    synchro.Now[tz.AsiaTokyo](),
		UpdatedAt:    synchro.Now[tz.AsiaTokyo](),
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
		Expiry:       token.Expiry.StdTime(),
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

type listLanguageResponce struct {
	Data       []domain.Language `json:"data"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}

func (c *client) ListLanguages(ctx context.Context) ([]domain.Language, error) {
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

	var languages listLanguageResponce
	err = json.Unmarshal(raw, &languages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return languages.Data, nil
}

type nowLanguageResponce struct {
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

func (c *client) NowLanguage(ctx context.Context) (string, error) {
	if c.wclient == nil {
		return "", fmt.Errorf("client is not set")
	}

	endpoint, err := url.Parse("https://wakatime.com/api/v1/users/current/durations")
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	query := endpoint.Query()
	query.Set("date", synchro.Now[tz.AsiaTokyo]().Format("2006-01-02"))
	query.Set("slice_by", "language")
	endpoint.RawQuery = query.Encode()

	resp, err := c.wclient.Get(endpoint.String())
	if err != nil {
		return "", fmt.Errorf("failed to get languages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get languages: %v", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var languages nowLanguageResponce
	err = json.Unmarshal(raw, &languages)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(languages.Data) == 0 {
		slog.Warn("no language")
		return "", nil
	}

	l := languages.Data[len(languages.Data)-1]
	lastTime := synchro.Unix[tz.AsiaTokyo](int64(math.Floor(l.Time+l.Duration)), 0)

	if lastTime.Before(synchro.Now[tz.AsiaTokyo]().Add(-10 * time.Minute)) {
		slog.Warn("last language is too old", "lastLanguage", l.Language, "lastTime", lastTime)
		return "", nil
	}
	return l.Language, nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/repository"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/service"
	"github.com/walnuts1018/wakatime-to-slack-profile/util"
	"golang.org/x/oauth2"
)

const (
	wakatimeAuthEndpoint  = "https://wakatime.com/oauth/authorize"
	wakatimeTokenEndpoint = "https://wakatime.com/oauth/token"
	wakatimeRedirectPath  = "/wakatime_callback"

	slackAuthEndpoint = "https://slack.com/oauth/authorize"
)

var (
	wakatimeScopes = []string{
		"read_stats",
		"read_logged_time",
		"email",
	}

	slackScopes     = []string{}
	slackUserScopes = []string{
		"users.profile:write",
		"emoji:read",
	}
)

type UpdateStatus struct {
	wakatimeUserRepo  repository.WakatimeUser
	slackUserRepo     repository.SlackUser
	userRepo          repository.User
	wakatimeTokenRepo repository.WakatimeToken
	slackTokenRepo    repository.SlackToken

	connectUserService service.ConnectUserService
	slackEmojiService  service.SlackEmojiService

	wakatimeOauth2Config *oauth2.Config

	wakatimeOauth2Clients map[string]*http.Client // wakatimeUserID: oauth2Client
	slackTokenCache       map[string]string       // slackUserID: slackToken

	slackAuth repository.SlackAuth

	noLanguageDuration time.Duration
}

func NewUpdateStatus(
	cfg config.Config,
	wakatimeUserRepo repository.WakatimeUser,
	slackUserRepo repository.SlackUser,
	userRepo repository.User,
	wakatimeTokenRepo repository.WakatimeToken,
	slackTokenRepo repository.SlackToken,
	connectUserService service.ConnectUserService,
	slackEmojiService service.SlackEmojiService,
	slackAuth repository.SlackAuth,
) *UpdateStatus {

	redirectURL, err := url.JoinPath(cfg.ServerURL, wakatimeRedirectPath)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to join url: %v", err), "server_url", cfg.ServerURL)
		os.Exit(1)
	}

	return &UpdateStatus{
		wakatimeUserRepo:   wakatimeUserRepo,
		slackUserRepo:      slackUserRepo,
		userRepo:           userRepo,
		wakatimeTokenRepo:  wakatimeTokenRepo,
		slackTokenRepo:     slackTokenRepo,
		connectUserService: connectUserService,
		slackEmojiService:  slackEmojiService,
		wakatimeOauth2Config: &oauth2.Config{
			ClientID:     cfg.WakatimeAppID,
			ClientSecret: cfg.WakatimeAppSecret,
			Endpoint:     oauth2.Endpoint{AuthURL: wakatimeAuthEndpoint, TokenURL: wakatimeTokenEndpoint},
			RedirectURL:  redirectURL,
			Scopes:       wakatimeScopes,
		},
		noLanguageDuration: cfg.NoLanguageDuration,
		slackAuth:          slackAuth,
	}
}

func (s *UpdateStatus) UpdateStatus(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}

	if user.WakatimeUserID == "" {
		return errors.New("wakatime user not found")
	}

	if user.SlackUserID == "" {
		return errors.New("slack user not found")
	}

	var client *http.Client
	if c, ok := s.wakatimeOauth2Clients[user.WakatimeUserID]; ok {
		client = c
	} else {
		client, err = s.GetOauthClient(ctx, user.WakatimeUserID)
		if err != nil {
			return fmt.Errorf("failed to get oauth client: %w", err)
		}
		s.wakatimeOauth2Clients[user.WakatimeUserID] = client
	}

	activity, err := s.wakatimeUserRepo.GetLastActivity(user.WakatimeUserID, client)
	if err != nil {
		return fmt.Errorf("failed to get last activity: %w", err)
	}
	var language string
	if synchro.Now[tz.AsiaTokyo]().Sub(activity.End) > s.noLanguageDuration {
		language = ""
	} else {
		language = activity.Language
	}

	var slackToken string
	cachedSlackToken, ok := s.slackTokenCache[user.SlackUserID]
	if ok {
		slackToken = cachedSlackToken
	} else {
		token, err := s.slackTokenRepo.GetSlackToken(user.SlackUserID)
		if err != nil {
			return fmt.Errorf("failed to get slack token: %w", err)
		}
		s.slackTokenCache[user.SlackUserID] = token.AccessToken
	}

	slackUser, err := s.slackUserRepo.GetUser(user.SlackUserID, slackToken)
	if err != nil {
		return fmt.Errorf("failed to get slack user: %w", err)
	}
	emoji := s.DetectEmoji(userID, slackUser.TeamID, language)

	if err := s.slackUserRepo.SetCustomStatus(slackUser.ID, slackToken, emoji, activity.Language); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	return nil
}

func (s *UpdateStatus) GetOauthClient(ctx context.Context, wakatimeUserID string) (*http.Client, error) {
	savedToken, err := s.wakatimeTokenRepo.GetWakatimeToken(wakatimeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wakatime token: %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  savedToken.AccessToken,
		RefreshToken: savedToken.RefreshToken,
		Expiry:       savedToken.ExpiresAt.StdTime(),
	}

	tokenSrc := s.wakatimeOauth2Config.TokenSource(ctx, token)
	myTokenSource := &tokenSource{
		src:               tokenSrc,
		wakatimeTokenRepo: s.wakatimeTokenRepo,
		wakatimeUserID:    wakatimeUserID,
	}

	reuseSrc := oauth2.ReuseTokenSource(token, myTokenSource)
	client := oauth2.NewClient(ctx, reuseSrc)

	return client, nil
}

func (s *UpdateStatus) DetectEmoji(userID, slackTeamID, language string) string {
	// カスタム絵文字
	overrides, err := s.userRepo.GetAllEmojis(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get custom emojis: %v", err), slog.String("user_id", userID))
	}
	emoji, ok := overrides[language]
	if ok {
		exists, err := s.slackEmojiService.IsEmojiExist(slackTeamID, emoji)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to check emoji: %v", err), slog.String("team_id", slackTeamID))
			exists = false
		}

		if exists {
			return emoji
		}
	}

	// sloth emoji
	if language == "" {
		// sloth emojiは環境によって表示が変わるので、代替絵文字の:namakemono:を使う
		exists, err := s.slackEmojiService.IsEmojiExist(slackTeamID, "namakemono")
		if err != nil {
			slog.Error(fmt.Sprintf("failed to check emoji: %v", err), slog.String("team_id", slackTeamID))
			exists = false
		}

		if exists {
			return "namakemono"
		} else {
			return "sloth"
		}
	}

	// そのまま
	exists, err := s.slackEmojiService.IsEmojiExist(slackTeamID, language)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check emoji: %v", err), slog.String("team_id", slackTeamID))
		exists = false
	}

	if exists {
		return language
	}

	// 全部小文字
	exists, err = s.slackEmojiService.IsEmojiExist(slackTeamID, strings.ToLower(language))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to check emoji: %v", err), slog.String("team_id", slackTeamID))
		exists = false
	}

	if exists {
		return strings.ToLower(language)
	}

	// 未知
	return "question"
}

func (s *UpdateStatus) WakatimeAuth() (string, string, error) {
	state, err := util.RandStr(64)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random string: %w", err)
	}

	redirect := s.wakatimeOauth2Config.AuthCodeURL(state)
	return state, redirect, nil
}

// return: wakatimeUserID, error
func (s *UpdateStatus) WakatimeCallback(ctx context.Context, code string) (string, error) {
	token, err := s.wakatimeOauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	cfg := model.WakatimeToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    synchro.In[tz.AsiaTokyo](token.Expiry),
		UpdatedAt:    synchro.Now[tz.AsiaTokyo](),
	}

	if err := s.wakatimeTokenRepo.UpdateWakatimeToken(cfg); err != nil {
		return "", fmt.Errorf("failed to update wakatime token: %w", err)
	}

	client := s.wakatimeOauth2Config.Client(ctx, token)
	user, err := s.wakatimeUserRepo.GetCurrentUser(client)
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	return user.ID, nil
}

func (s *UpdateStatus) SlackAuth() string {
	url, err := url.Parse(slackAuthEndpoint)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse url: %v", err), "slack_auth_endpoint", slackAuthEndpoint)
		os.Exit(1)
	}

	query := url.Query()
	query.Set("client_id", s.wakatimeOauth2Config.ClientID)
	query.Set("scope", strings.Join(slackScopes, ","))
	query.Set("user_scope", strings.Join(slackUserScopes, ","))

	url.RawQuery = query.Encode()

	return url.String()
}

// return: slackUserID, error
func (s *UpdateStatus) SlackCallback(ctx context.Context, code string) (string, error) {
	token, userID, err := s.slackAuth.GetOAuthV2Response(code)
	if err != nil {
		return "", fmt.Errorf("failed to slack callback: %w", err)
	}

	slackToken := model.SlackToken{
		SlackUserID: userID,
		AccessToken: token,
	}

	if err := s.slackTokenRepo.UpdateSlackToken(slackToken); err != nil {
		return "", fmt.Errorf("failed to update slack token: %w", err)
	}

	return userID, nil
}

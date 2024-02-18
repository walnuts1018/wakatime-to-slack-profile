package wakatime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/domain/model"
)

const base = "https://wakatime.com/api/v1/"

type Client struct {
	config config.Config
}

func NewClient(cfg config.Config) (Client, error) {
	return Client{config: cfg}, nil
}

func (c Client) GetUser(wakatimeUserID string, oauth2Client *http.Client) (model.WakatimeUser, error) {
	if oauth2Client == nil {
		return model.WakatimeUser{}, fmt.Errorf("oauth2 client is nil")
	}

	url, err := url.JoinPath(base, "users", wakatimeUserID)
	if err != nil {
		return model.WakatimeUser{}, fmt.Errorf("failed to join url: %w", err)
	}

	resp, err := oauth2Client.Get(url)
	if err != nil {
		return model.WakatimeUser{}, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	var user model.WakatimeUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return model.WakatimeUser{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return user, nil
}

func (c Client) GetCurrentUser(oauth2Client *http.Client) (model.WakatimeUser, error) {
	return c.GetUser("current", oauth2Client)
}

func (c Client) GetLastActivity(wakatimeUserID string, oauth2Client *http.Client) (model.WakatimeLastActivity, error) {
	if oauth2Client == nil {
		return model.WakatimeLastActivity{}, fmt.Errorf("oauth2 client is nil")
	}

	urlStr, err := url.JoinPath(base, "users", wakatimeUserID, "durations")
	if err != nil {
		return model.WakatimeLastActivity{}, fmt.Errorf("failed to join url: %w", err)
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return model.WakatimeLastActivity{}, fmt.Errorf("failed to parse url: %w", err)
	}

	query := url.Query()
	query.Set("date", synchro.Now[tz.AsiaTokyo]().Format("2006-01-02"))
	query.Set("slice_by", "language")
	url.RawQuery = query.Encode()

	resp, err := oauth2Client.Get(url.String())
	if err != nil {
		return model.WakatimeLastActivity{}, fmt.Errorf("failed to get last activity: %w", err)
	}
	defer resp.Body.Close()

	var durations durationsResponce
	if err := json.NewDecoder(resp.Body).Decode(&durations); err != nil {
		return model.WakatimeLastActivity{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(durations.Data) == 0 {
		return model.WakatimeLastActivity{}, nil
	}

	lastData := durations.Data[len(durations.Data)-1]

	startTime := synchro.Unix[tz.AsiaTokyo](int64(lastData.Time), 0)
	endTime := synchro.Unix[tz.AsiaTokyo](int64(lastData.Time+lastData.Duration), 0)

	return model.WakatimeLastActivity{
		Language: lastData.Language,
		Project:  lastData.Project,
		Start:    startTime,
		End:      endTime,
	}, nil
}

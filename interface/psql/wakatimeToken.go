package psql

import (
	"fmt"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

func (c Client) UpdateWakatimeToken(wakatimeUserID, accessToken, refreshToken string, expiresAt synchro.Time[tz.AsiaTokyo]) error {
	token := WakatimeToken{
		WakatimeUserID: wakatimeUserID,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		ExpiresAt:      expiresAt,
	}
	result := c.db.Create(&token)
	if result.Error != nil {
		return fmt.Errorf("failed to update wakatime token: %w", result.Error)
	}
	return nil
}

func (c Client) GetWakatimeToken(wakatimeUserID string) (WakatimeToken, error) {
	var token WakatimeToken
	result := c.db.Where("wakatime_user_id = ?", wakatimeUserID).First(&token)
	if result.Error != nil {
		return WakatimeToken{}, fmt.Errorf("failed to get wakatime token: %w", result.Error)
	}
	return token, nil
}

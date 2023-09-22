package usecase

import (
	"fmt"
	"log/slog"
	"strings"
)

func (u *Usecase) SetUserCustomStatus(language string) error {
	if language == "" {
		err := u.slackClient.SetUserCustomStatus("sloth")
		if err == nil {
			slog.Info("set user custom status", "emoji", "🦥")
			return nil
		}
	}

	override, ok := u.emojiOverides[language]
	if ok {
		err := u.slackClient.SetUserCustomStatus(override)
		if err == nil {
			slog.Info("set user custom status", "emoji", override)
			return nil
		}
	}

	err := u.slackClient.SetUserCustomStatus(language)
	if err == nil {
		slog.Info("set user custom status", "emoji", language)
		return nil
	}
	err = u.slackClient.SetUserCustomStatus(strings.ToLower(language))
	if err == nil {
		slog.Info("set user custom status", "emoji", language)
		return nil
	}

	err = u.slackClient.SetUserCustomStatus("question")
	if err == nil {
		slog.Info("set user custom status", "emoji", "❓")
		return nil
	}

	return fmt.Errorf("failed to find emoji: %v", language)
}

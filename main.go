package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/handler"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/psql"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/slack"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/wakatime"
	"github.com/walnuts1018/wakatime-to-slack-profile/usecase"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	ctx, canecel := context.WithCancel(context.Background())
	defer canecel()

	psqClient, err := psql.NewClient()
	if err != nil {
		slog.Error("Error creating psql client", "error", err)
		os.Exit(1)
	}
	defer psqClient.Close()

	wakatimeClient := wakatime.NewOauth2Client()

	slackClient := slack.NewClient()

	usecase := usecase.NewUsecase(wakatimeClient, psqClient, slackClient)
	err = usecase.SetToken(ctx)
	if err != nil {
		slog.Warn("failed to set token", "error", err)
	}

	go func() {
		err := usecase.SetLanguage(ctx)
		if err != nil {
			slog.Error("Failed to set language", "error", err)
			os.Exit(1)
		}

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				slog.Info("context done")
				os.Exit(0)
			case <-ticker.C:
				err := usecase.SetLanguage(ctx)
				if err != nil {
					slog.Error("Failed to set language", "error", err)
					os.Exit(1)
				}
			}
		}
	}()

	handler, err := handler.NewHandler(usecase)
	if err != nil {
		slog.Error("Error loading handler: %v", "error", err)
		os.Exit(1)
	}
	err = handler.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
	if err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/handler"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/psql"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/slack"
	"github.com/walnuts1018/wakatime-to-slack-profile/infra/wakatime"
	"github.com/walnuts1018/wakatime-to-slack-profile/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
		Level:      cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	ctx, canecel := context.WithCancel(context.Background())
	defer canecel()

	psqClient, err := psql.NewClient(cfg)
	if err != nil {
		slog.Error("Error creating psql client", "error", err)
		os.Exit(1)
	}
	defer psqClient.Close()

	wakatimeClient, err := wakatime.NewOauth2Client(cfg)
	if err != nil {
		slog.Error("Error creating wakatime client", "error", err)
		os.Exit(1)
	}

	slackClient := slack.NewClient(cfg)

	usecase := usecase.NewUsecase(wakatimeClient, psqClient, slackClient, map[string]string{})
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

	handler, err := handler.NewHandler(cfg, usecase, logger)
	if err != nil {
		slog.Error("Error loading handler: %v", "error", err)
		os.Exit(1)
	}
	if err := handler.Run(fmt.Sprintf(":%v", cfg.ServerPort)); err != nil {
		slog.Error("failed to run handler", "error", err)
		os.Exit(1)
	}
}

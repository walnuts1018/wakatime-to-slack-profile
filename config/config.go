package config

import (
	"crypto/rand"
	"flag"
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	WakatimeAppID     string `env:"WAKATIME_APP_ID,required"`
	WakatimeAppSecret string `env:"WAKATIME_CLIENT_SECRET,required"`
	CookieSecret      string `env:"COOKIE_SECRET"`

	SlackClientID     string `env:"SLACK_CLIENT_ID,required"`
	SlackClientSecret string `env:"SLACK_CLIENT_SECRET,required"`

	PSQLEndpoint string `env:"PSQL_ENDPOINT" envDefault:"localhost"`
	PSQLPort     string `env:"PSQL_PORT" envDefault:"5432"`
	PSQLDatabase string `env:"PSQL_DATABASE" envDefault:"wakatime_to_slack"`
	PSQLUser     string `env:"PSQL_USER" envDefault:"postgres"`
	PSQLPassword string `env:"PSQL_PASSWORD" envDefault:"postgres"`
	PSQLSSLMode  string `env:"PSQL_SSL_MODE" envDefault:"disable"`

	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	ServerURL  string `env:"SERVER_URL" envDefault:"https://wakatime.walnuts.dev"`

	LogLevelString string `env:"LOG_LEVEL" envDefault:"info"`
	LogLevel       slog.Level

	EmojiCacheDuration time.Duration `env:"EMOJI_CACHE_DURATION" envDefault:"24h"`
	NoLanguageDuration time.Duration `env:"NO_LANGUAGE_DURATION" envDefault:"10m"` // コードを書いていないと判断される時間
}

func Load() (Config, error) {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()

	cfg := Config{}
	if serverport != nil {
		cfg.ServerPort = *serverport
	}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	// parse log level
	switch strings.ToLower(cfg.LogLevelString) {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "info":
		cfg.LogLevel = slog.LevelInfo
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		slog.Warn("Invalid log level, use default level: info")
		cfg.LogLevel = slog.LevelInfo
	}

	// set CookieSecret
	if cfg.CookieSecret == "" {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			return cfg, err
		}

		var result string
		for _, v := range b {
			result += string(letters[int(v)%len(letters)])
		}
		cfg.CookieSecret = result
	}

	return cfg, nil
}

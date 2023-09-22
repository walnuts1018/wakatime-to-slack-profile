package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config_t struct {
	WakatimeAppID     string `env:"WAKATIME_APP_ID"`
	WakatimeAppSecret string `env:"WAKATIME_CLIENT_SECRET"`
	CookieSecret      string `env:"COOKIE_SECRET"`

	PSQLEndpoint string `env:"PSQL_ENDPOINT"`
	PSQLPort     string `env:"PSQL_PORT"`
	PSQLDatabase string `env:"PSQL_DATABASE"`
	PSQLUser     string `env:"PSQL_USER"`
	PSQLPassword string `env:"PSQL_PASSWORD"`

	SlackAccessToken string `env:"SLACK_ACCESS_TOKEN"`

	ServerPort string
}

var Config = Config_t{}

func LoadConfig() error {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()
	Config.ServerPort = *serverport

	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	t := reflect.TypeOf(Config)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		tag, ok := t.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		v, ok := os.LookupEnv(tag)
		if !ok {
			return fmt.Errorf("%s is not set", tag)
		}
		reflect.ValueOf(&Config).Elem().FieldByName(fieldName).SetString(v)
	}
	return nil
}

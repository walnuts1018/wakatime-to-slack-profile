package psql

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg config.Config) (Client, error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", cfg.PSQLEndpoint, cfg.PSQLPort, cfg.PSQLUser, cfg.PSQLPassword, cfg.PSQLDatabase, cfg.PSQLSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return Client{}, fmt.Errorf("failed to open database: %w", err)
	}

	return Client{db: db}, nil
}

func (c Client) Init() error {
	if err := c.db.AutoMigrate(&CustomEmojis{}, &User{}, &WakatimeToken{}, &SlackToken{}, &SlackEmoji{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

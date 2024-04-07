package main

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/pkarpovich/tg-link-keeper-bot/app/bot/linkstore"
	"github.com/pkarpovich/tg-link-keeper-bot/app/events"
	"log"
	"os"
)

type Config struct {
	Telegram struct {
		Token      string  `env:"TELEGRAM_TOKEN"`
		SuperUsers []int64 `env:"TELEGRAM_SUPER_USERS" envSeparator:","`
	}
	LinkStore struct {
		Url string `env:"LINK_STORE_URL"`
	}
}

func main() {
	log.Printf("[INFO] starting app...")

	config, err := prepareConfig()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	if err := execute(config); err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

func prepareConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &cfg, nil
}

func execute(config *Config) error {
	tbAPI, err := tbapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return fmt.Errorf("failed to create Telegram events: %w", err)
	}

	linkdingClient := linkstore.NewLinkStoreClient(config.LinkStore.Url)

	tgListener := &events.TelegramListener{
		SuperUsers: config.Telegram.SuperUsers,
		TbAPI:      tbAPI,
		Bot:        linkdingClient,
	}

	if err := tgListener.Do(); err != nil {
		return fmt.Errorf("failed to start Telegram listener: %w", err)
	}

	return nil
}

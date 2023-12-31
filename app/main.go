package main

import (
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golobby/dotenv"
	"github.com/pkarpovich/tg-link-keeper-bot/app/bot/linkstore"
	"github.com/pkarpovich/tg-link-keeper-bot/app/events"
	"log"
	"os"
)

var config struct {
	Telegram struct {
		Token string `env:"TELEGRAM_TOKEN"`
	}
	LinkStore struct {
		Url string `env:"LINK_STORE_URL"`
	}
}

func main() {
	log.Printf("[INFO] starting app...")

	if err := prepareConfig(); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	if err := execute(); err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

func prepareConfig() error {
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}

	err = dotenv.NewDecoder(file).Decode(&config)
	if err != nil {
		return fmt.Errorf("failed to decode .env file: %w", err)
	}

	return nil
}

func execute() error {
	tbAPI, err := tbapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return fmt.Errorf("failed to create Telegram events: %w", err)
	}

	linkdingClient := linkstore.NewLinkStoreClient(config.LinkStore.Url)

	tgListener := &events.TelegramListener{
		TbAPI: tbAPI,
		Bot:   linkdingClient,
	}

	if err := tgListener.Do(); err != nil {
		return fmt.Errorf("failed to start Telegram listener: %w", err)
	}

	return nil
}

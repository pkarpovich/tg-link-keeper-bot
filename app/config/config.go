package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Telegram struct {
		Token      string  `env:"TELEGRAM_TOKEN"`
		SuperUsers []int64 `env:"TELEGRAM_SUPER_USERS" envSeparator:","`
	}
	LinkStore struct {
		Url     string `env:"LINK_STORE_URL"`
		DryMode bool   `env:"LINK_STORE_DRY_MODE" envDefault:"false"`
	}
}

func Init() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[WARN] error while loading .env file: %v", err)
	}

	var cfg Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

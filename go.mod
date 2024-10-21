module github.com/pkarpovich/tg-link-keeper-bot

go 1.23

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jimmysawczuk/recon v0.0.0-20240723135856-0ca09c7808a6
	github.com/joho/godotenv v1.5.1
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api/v5 v5.0.0-20241013102643-36756d99d4ae

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

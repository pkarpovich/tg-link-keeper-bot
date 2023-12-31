module github.com/pkarpovich/tg-link-keeper-bot

go 1.21

require (
	github.com/caarlos0/env/v10 v10.0.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jimmysawczuk/recon v0.0.0-20230225193537-3366d8a9e56f
	github.com/joho/godotenv v1.5.1
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/OvyFlash/telegram-bot-api/v5 v5.0.0-20231230151827-6d16deaa376e

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.17.0 // indirect
)

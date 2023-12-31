# Telegram Link Keeper Bot

Link Keeper Bot is a Telegram bot designed to save links or memo files in Cubox. It processes different types of
messages including URLs and forwarded messages.

## Features

- Save URLs and text memos.
- Process forwarded messages from Telegram.
- Restrict access to specified 'Super Users'.
- Integrate with an external link storage service (Cubox).

## Requirements

- Docker and Docker Compose
- Telegram Bot API token
- Access to a Cubox instance or similar service for link storage

## Installation

1. Clone the repository:
```bash
git clone git@github.com:pkarpovich/tg-link-keeper-bot.git
```
2. Navigate to the project directory:
```bash
cd tg-link-keeper-bot
```

## Configuration
1. Create a .env file in the project root with the following contents:
```bash
TELEGRAM_TOKEN=your_telegram_bot_token
TELEGRAM_SUPER_USERS=user_id1,user_id2
LINK_STORE_URL=your_link_store_endpoint
```
2. Adjust the values for `TELEGRAM_TOKEN`, `TELEGRAM_SUPER_USERS`, and `LINK_STORE_URL` according to your setup.

## Docker Deployment
1. Use the provided Docker Compose file to deploy the bot:
```bash
docker compose pull
docker compose up -d
```
2. The bot will be running in a container named tg-link-keeper-bot.

## Usage
- Send a URL or a forwarded message to the bot to save it in Cubox.
- Send a text message to the bot to save it as a memo in Cubox.
- Interact with the bot in Telegram. Use the /ping command to check its status.

## Commands
- `/ping` - Check if the bot is online.

## License
[MIT](https://choosealicense.com/licenses/mit/)
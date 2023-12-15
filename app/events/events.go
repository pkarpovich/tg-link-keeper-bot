package events

import (
	"encoding/json"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkarpovich/tg-linkding/app/bot"
	"log"
)

type Bot interface {
	OnMessage(msg bot.Message) error
}

type TbAPI interface {
	GetUpdatesChan(config tbapi.UpdateConfig) tbapi.UpdatesChannel
}

type TelegramListener struct {
	TbAPI TbAPI
	Bot   Bot
}

func (tl *TelegramListener) Do() error {
	u := tbapi.NewUpdate(0)
	u.Timeout = 60

	updates := tl.TbAPI.GetUpdatesChan(u)

	for {
		select {

		case update, ok := <-updates:
			if !ok {
				return fmt.Errorf("telegram update chan closed")
			}

			if update.Message == nil {
				continue
			}

			if err := tl.processEvent(update); err != nil {
				return fmt.Errorf("failed to process event: %w", err)
			}
		}
	}
}

func (tl *TelegramListener) processEvent(update tbapi.Update) error {
	msgJSON, errJSON := json.Marshal(update.Message)
	if errJSON != nil {
		return fmt.Errorf("failed to marshal update.Message to json: %w", errJSON)
	}
	log.Printf("[DEBUG] %s", string(msgJSON))

	msg := tl.transform(update.Message)
	if err := tl.Bot.OnMessage(msg); err != nil {
		return fmt.Errorf("failed to process message: %w", err)
	}

	return nil
}

func (tl *TelegramListener) transform(message *tbapi.Message) bot.Message {
	return bot.Message{
		ID:     message.MessageID,
		From:   bot.User{},
		ChatID: message.Chat.ID,
		HTML:   message.Text,
		Text:   message.Text,
		Sent:   message.Time(),
	}
}

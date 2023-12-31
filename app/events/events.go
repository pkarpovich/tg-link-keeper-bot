package events

import (
	"encoding/json"
	"errors"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkarpovich/tg-link-keeper-bot/app/bot"
	"log"
)

const (
	PingCommand = "ping"
)

type Bot interface {
	OnMessage(msg bot.Message) error
}

type TbAPI interface {
	GetUpdatesChan(config tbapi.UpdateConfig) tbapi.UpdatesChannel
	Send(c tbapi.Chattable) (tbapi.Message, error)
}

type TelegramListener struct {
	SuperUsers []int64
	TbAPI      TbAPI
	Bot        Bot
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
				log.Printf("[ERROR] %v", err)
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

	if !tl.isSuperUser(update.Message.From.ID) {
		log.Printf("[DEBUG] user %d is not super user", update.Message.From.ID)

		msg := tbapi.NewMessage(update.Message.Chat.ID, "I don't know you 🤷‍")
		_, err := tl.TbAPI.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		return nil
	}

	switch update.Message.Command() {
	case PingCommand:
		tl.handlePingCommand(update)
		return nil
	}

	msg := tl.transform(update.Message)
	if err := tl.Bot.OnMessage(msg); err != nil {
		errMsg := tbapi.NewMessage(update.Message.Chat.ID, "💥 Error: "+err.Error())
		_, err := tl.TbAPI.Send(errMsg)
		if err != nil {
			return fmt.Errorf("failed to send error message: %w", err)
		}

		return errors.New(errMsg.Text)
	}

	resultMsg := tbapi.NewMessage(update.Message.Chat.ID, "💾 Saved!")
	_, err := tl.TbAPI.Send(resultMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (tl *TelegramListener) transform(message *tbapi.Message) bot.Message {
	msg := bot.Message{
		ID:     message.MessageID,
		From:   bot.User{},
		ChatID: message.Chat.ID,
		HTML:   message.Text,
		Text:   message.Text,
		Sent:   message.Time(),
	}

	if message.ForwardFromChat != nil {
		msg.ForwardFromMessageID = message.ForwardFromMessageID
		msg.ForwardFromChat = &bot.Chat{
			ID:       message.ForwardFromChat.ID,
			Title:    message.ForwardFromChat.Title,
			UserName: message.ForwardFromChat.UserName,
		}
	}

	return msg
}

func (tl *TelegramListener) handlePingCommand(update tbapi.Update) {
	msg := tbapi.NewMessage(update.Message.Chat.ID, "🏓 Pong!")
	_, err := tl.TbAPI.Send(msg)
	if err != nil {
		log.Printf("[ERROR] failed to send message: %v", err)
	}
}

func (tl *TelegramListener) isSuperUser(userID int64) bool {
	for _, su := range tl.SuperUsers {
		if su == userID {
			return true
		}
	}

	return false
}

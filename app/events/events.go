package events

import (
	"encoding/json"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkarpovich/tg-link-keeper-bot/app/bot"
	"iter"
	"log"
	"time"
)

const (
	TickerTimeout = 2 * time.Second
	PingCommand   = "ping"
)

type Bot interface {
	OnMessage(msg bot.Message) iter.Seq[bot.Response]
}

type TbAPI interface {
	GetUpdatesChan(config tbapi.UpdateConfig) tbapi.UpdatesChannel
	Send(c tbapi.Chattable) (tbapi.Message, error)
	Request(c tbapi.Chattable) (*tbapi.APIResponse, error)
}

type TelegramListener struct {
	SuperUsers []int64
	TbAPI      TbAPI
	Bot        bot.MultiBot
}

func (tl *TelegramListener) Do() error {
	u := tbapi.NewUpdate(0)
	u.Timeout = 60

	updates := tl.TbAPI.GetUpdatesChan(u)
	batchMap := make(map[int64][]tbapi.Update)
	ticker := time.NewTicker(TickerTimeout)

	for {
		select {

		case update, ok := <-updates:
			if !ok {
				return fmt.Errorf("telegram update chan closed")
			}

			if update.Message == nil {
				continue
			}

			userID := update.Message.From.ID
			if !tl.isSuperUser(userID) {
				if err := tl.handleNonSuperUser(update); err != nil {
					return err
				}

				return nil
			}

			switch update.Message.Command() {
			case PingCommand:
				return tl.handlePingCommand(update)
			}

			batchMap[userID] = append(batchMap[userID], update)
		case <-ticker.C:
			for userID, batch := range batchMap {
				go tl.processBatch(batch)
				delete(batchMap, userID)
			}
		}
	}
}

func (tl *TelegramListener) processBatch(batch []tbapi.Update) {
	if messagesBelongToSameMediaGroup(batch) {
		if err := tl.processEvent(getMainMessageFromMediaGroup(batch)); err != nil {
			log.Printf("[ERROR] %v", err)
		}

		return
	}

	for _, update := range batch {
		if err := tl.processEvent(update); err != nil {
			log.Printf("[ERROR] %v", err)
		}
	}
}

func (tl *TelegramListener) processEvent(update tbapi.Update) error {
	msgJSON, errJSON := json.Marshal(update.Message)
	if errJSON != nil {
		return fmt.Errorf("failed to marshal update.Message to json: %w", errJSON)
	}
	log.Printf("[DEBUG] %s", string(msgJSON))

	for resp := range tl.Bot.OnMessage(tl.transform(update.Message)) {
		if resp.Reaction != nil {
			reactionMsg := tbapi.SetMessageReactionConfig{
				BaseChatMessage: tbapi.BaseChatMessage{
					ChatConfig: tbapi.ChatConfig{
						ChatID: resp.ChatID,
					},
					MessageID: resp.Reaction.MessageID,
				},
				Reaction: []tbapi.ReactionType{
					{
						Type:  "emoji",
						Emoji: resp.Reaction.Emoji,
					},
				},
				IsBig: false,
			}

			_, err := tl.TbAPI.Request(reactionMsg)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}

		if resp.Text != "" {
			msg := tbapi.NewMessage(resp.ChatID, resp.Text)
			_, err := tl.TbAPI.Send(msg)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
		}
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

	if len(message.Caption) > 0 {
		msg.Text = message.Caption
	}

	if message.ForwardOrigin != nil {
		origin := message.ForwardOrigin

		switch message.ForwardOrigin.Type {
		case tbapi.MessageOriginChannel:
			msg.Url = fmt.Sprintf("https://t.me/%s/%d", origin.Chat.UserName, origin.MessageID)
		case tbapi.MessageOriginUser:
			msg.Text = fmt.Sprintf(
				"%s %s (%s):\n%s",
				origin.SenderUser.FirstName,
				origin.SenderUser.LastName,
				origin.SenderUser.UserName,
				message.Text,
			)
		case tbapi.MessageOriginHiddenUser:
			msg.Text = fmt.Sprintf("%s:\n%s", origin.SenderUserName, message.Text)
		}
	}

	return msg
}

func (tl *TelegramListener) handlePingCommand(update tbapi.Update) error {
	msg := tbapi.NewMessage(update.Message.Chat.ID, "🏓 Pong!")
	_, err := tl.TbAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (tl *TelegramListener) handleNonSuperUser(update tbapi.Update) error {
	log.Printf("[DEBUG] user %d is not super user", update.Message.From.ID)

	msg := tbapi.NewMessage(update.Message.Chat.ID, "I don't know you 🤷‍")
	_, err := tl.TbAPI.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (tl *TelegramListener) isSuperUser(userID int64) bool {
	for _, su := range tl.SuperUsers {
		if su == userID {
			return true
		}
	}

	return false
}

func messagesBelongToSameMediaGroup(updates []tbapi.Update) bool {
	mediaGroupID := ""

	for _, upd := range updates {
		if mediaGroupID == "" {
			mediaGroupID = upd.Message.MediaGroupID
		}

		if upd.Message.MediaGroupID != mediaGroupID {
			return false
		}
	}

	return true
}

func getMainMessageFromMediaGroup(updates []tbapi.Update) tbapi.Update {
	var mainUpdate tbapi.Update

	for _, upd := range updates {
		if len(upd.Message.Photo) > 0 && len(upd.Message.Caption) > 0 {
			mainUpdate = upd
		}
	}

	return mainUpdate
}

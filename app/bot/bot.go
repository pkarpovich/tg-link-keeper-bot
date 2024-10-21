package bot

import (
	"iter"
	"time"
)

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"user_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type Message struct {
	ID     int
	From   User
	ChatID int64
	Sent   time.Time
	HTML   string `json:",omitempty"`
	Text   string `json:",omitempty"`
	Url    string
}

type ResponseReaction struct {
	Emoji     string
	MessageID int
}

type Response struct {
	Reaction *ResponseReaction
	ChatID   int64
	Text     string
}

type Bot interface {
	ShouldHandle(msg Message) bool
	OnMessage(msg Message) Response
}

type MultiBot []Bot

func (mb *MultiBot) OnMessage(msg Message) iter.Seq[Response] {
	return func(yield func(Response) bool) {
		for _, bot := range *mb {
			if !bot.ShouldHandle(msg) {
				continue
			}

			if !yield(bot.OnMessage(msg)) {
				return
			}
		}
	}
}

package bot

import (
	"fmt"
	"github.com/pkarpovich/tg-link-keeper-bot/app/providers"
	"log"
)

type LinkStore struct {
	cubox *providers.Cubox
}

func NewLinkstore(cubox *providers.Cubox) *LinkStore {
	return &LinkStore{
		cubox: cubox,
	}
}

func (l *LinkStore) ShouldHandle(msg Message) bool {
	return true
}

func (l *LinkStore) OnMessage(msg Message) Response {
	content := l.cubox.PrepareContent(msg)
	if content == nil {
		log.Printf("[DEBUG] empty content")
		return Response{}
	}

	if err := l.cubox.SaveLink(*content); err != nil {
		return Response{
			Text:   fmt.Sprintf("failed to save link: %v", err),
			ChatID: msg.ChatID,
		}
	}

	return Response{
		Reaction: &ResponseReaction{
			Emoji:     "👍",
			MessageID: msg.ID,
		},
		ChatID: msg.ChatID,
	}
}

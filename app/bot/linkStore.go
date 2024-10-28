package bot

import (
	"errors"
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
	content := l.cubox.PrepareContent(msg.Text, msg.Url)
	if content == nil {
		log.Printf("[DEBUG] empty content")
		return Response{}
	}

	if err := l.cubox.SaveLink(*content); err != nil {
		if errors.Is(err, providers.ErrDuplicatedLink) {
			return Response{
				Reaction: &ResponseReaction{
					MessageID: msg.ID,
					Emoji:     "👀",
				},
				ChatID: msg.ChatID,
			}
		}

		return Response{
			Text:   fmt.Sprintf("failed to save link: %v", err),
			ChatID: msg.ChatID,
		}
	}

	return Response{
		Reaction: &ResponseReaction{
			MessageID: msg.ID,
			Emoji:     "👍",
		},
		ChatID: msg.ChatID,
	}
}

package bot

import "time"

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"user_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Title    string `json:"title,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

type Message struct {
	ID                   int
	From                 User
	ChatID               int64
	Sent                 time.Time
	HTML                 string `json:",omitempty"`
	Text                 string `json:",omitempty"`
	ForwardFromChat      *Chat
	ForwardFromMessageID int
}

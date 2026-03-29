package models

type ChatMsg struct {
	UserID   int64  `json:"-"`
	FriendID int64  `json:"-"`
	ChatID   int64  `json:"chat_id"`
	Context  string `json:"context"`
}

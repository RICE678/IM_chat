package models

import "time"

type ChatMsg struct {
	UserID     int64     `json:"-"`
	ReceiverID int64     `json:"receiver_id" form:"receiver_id" binding:"required"`
	ChatID     int64     `json:"chat_id" form:"chat_id" binding:"required"`
	Context    string    `json:"context" form:"context" binding:"required"`
	MsgType    string    `json:"msg_type" form:"msg_type" binding:"required"`
	CreateTime time.Time `json:"-"`
}

type WsMsg struct {
	Type       string `json:"type"`
	Msg        string `json:"msg"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Timestamp  int64  `json:"timestamp"`
	MsgType    int    `json:"msy_type"`
}

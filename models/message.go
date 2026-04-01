package models

import "time"

type ChatMsg struct {
	ID         int64     `json:"-"`
	UserID     int64     `json:"-"`
	ReceiverID int64     `json:"receiver_id" form:"receiver_id" binding:"required"`
	CreateTime time.Time `json:"-"`
	SendStatus int       `json:"-"`
	MsgType    int       `json:"msg_type" form:"msg_type" binding:"required"`
	Context    string    `json:"context" form:"context" binding:"required"`
	IsDelete   int       `json:"-"`
	ReadStatus int       `json:"-"`
}

type WsMsg struct {
	Type       string `json:"type"`
	Msg        string `json:"msg"`
	SenderID   int64  `json:"-"`
	ReceiverID int64  `json:"receiver_id"`
	Timestamp  int64  `json:"-"`
	MsgType    int    `json:"msg_type"`
}

type HistoryMsg struct {
	Msg        []*ChatMsg `json:"msg"`
	SenderID   int64      `json:"-"`
	ReceiverID int64      `json:"receiver_id" form:"receiver_id" binding:"required"`
	Page       int        `json:"page" form:"page"`
	Size       int        `json:"size" form:"size"`
}
type HistoryResponse struct {
	Msg *HistoryMsg `json:"msg"`
	Err string      `json:"err"`
}

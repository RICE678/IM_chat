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
	// Type 消息类型，常用值：private / text / ack
	Type  string `json:"type"`
	MsgID int64  `json:"-"`
	// Msg 文本消息内容；文件消息时可传 file_url
	Msg string `json:"msg"`
	// FileURL 文件消息地址（msg_type=2 时建议必传）
	FileURL string `json:"file_url,omitempty"`
	// FileName 文件名（msg_type=2 时建议传）
	FileName   string `json:"file_name,omitempty"`
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Timestamp  int64  `json:"-"`
	// MsgType 1=文本消息，2=文件消息
	MsgType int `json:"msg_type"`
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

package models

import "time"

type ChatMsg struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id" form:"receiver_id" binding:"required"`
	CreateTime time.Time `json:"-"`
	Timestamp  int64     `json:"timestamp"`
	SendStatus int       `json:"-"`
	MsgType    int       `json:"msg_type" form:"msg_type" binding:"required"`
	// Keep both fields for compatibility: websocket uses "msg", old history uses "context".
	Msg        string `json:"msg,omitempty"`
	Context    string `json:"context,omitempty" form:"context"`
	IsDelete   int    `json:"-"`
	ReadStatus int    `json:"-"`
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
	Timestamp  int64  `json:"timestamp"`
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

type DelResponse struct {
	FriendID int64 `json:"friend_id"`
	UserID   int64
}

type ChatFileUploadResponse struct {
	Err string `json:"err"`
	// ID 与磁盘文件名中的雪花 id 一致，对应 docs.id
	ID int64 `json:"id,omitempty"`

	URL string `json:"url"`

	FileURL string `json:"file_url"`

	AbsURL   string `json:"abs_url,omitempty"`
	FileName string `json:"file_name"`
	MsgType  int    `json:"msg_type"`
}

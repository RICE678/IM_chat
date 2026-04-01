package models

import "time"

type MainFriend struct {
	FriendID    int64     `json:"friend_id"`
	LastMsgTime time.Time `json:"last_msg_time"`
	Unread      int       `json:"unread"`
}
type ContactResponse struct {
	Friend []MainFriend `json:"friend"`
	Err    string       `json:"error"`
}

type ReadContact struct {
	UserID   int64 `json:"-"`
	FriendID int64 `json:"friend_id" binding:"required"`
}

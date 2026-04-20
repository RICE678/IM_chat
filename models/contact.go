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

type ListContact struct {
	UserID        int64  `json:"user_id"`
	FriendID      int64  `json:"friend_id"`
	FriendPicture string `json:"friend_picture"`
	FriendName    string `json:"friend_name"`
	Remark        string `json:"remark,omitempty"`
}

type Contact struct {
	List []ListContact `json:"list"`
	Err  string        `json:"error"`
}

type FriendMain struct {
	UserID   int64 `json:"-"`
	FriendID int64 `json:"friend_id"`
}
type ContactMain struct {
	FriendID      int64  `json:"friend_id"`
	Remark        string `json:"remark,omitempty"`
	FriendEmail   string `json:"friend_email"`
	FriendPicture string `json:"friend-picture"`
	FriendName    string `json:"friend_name"`
	Gender        int    `json:"gender"`
	Signature     string `json:"signature,omitempty"`
}

type FriendRemark struct {
	UserID   int64  `json:"-"`
	FriendID int64  `json:"friend_id"`
	Remark   string `json:"remark,omitempty"`
}

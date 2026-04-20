package models

import "time"

type AppliSearch struct {
	ApplyID   int64  `json:"-"`
	UserID    int64  `json:"-"`
	SendEmail string `json:"send_email" binding:"required" form:"send_email"`
	SendID    int64  `json:"-"`
	Msg       string `json:"msg" binding:"required" form:"msg"`
	Status    int    `json:"-"`
}

type FindPerson struct {
	SendEmail string `json:"email" form:"email" `
	UserID    int64  `json:"-"`
	SendID    int64  `json:"-"`
}
type FindMiddle struct {
	SendEmail    string `json:"email"`
	SendName     string `json:"name"`
	SendID       int64  `json:"ID"`
	SendPictures string `json:"pictures"`
}

type FindEnd struct {
	Find *FindMiddle `json:"find,omitempty"`
	Err  string      `json:"error"`
}

type FindNamePerson struct {
	SendName  string `json:"name" form:"name" binding:"required"`
	UserID    int64  `json:"-"`
	SendID    int64  `json:"-"`
	SendEmail string `json:"-"`
}
type FindNameMiddle struct {
	SendEmail    string `json:"email"`
	SendName     string `json:"name"`
	SendID       int64  `json:"ID"`
	SendPictures string `json:"pictures"`
}

type FindNameEnd struct {
	Find []FindNameMiddle `json:"find,omitempty"`
	Err  string           `json:"error"`
}

type Apply struct {
	FromEmail   string    `json:"from_email"`
	FromID      int64     `json:"from_id"`
	SendEmail   string    `json:"send_email"`
	Msg         string    `json:"msg"`
	Time        time.Time `json:"time"`
	SendName    string    `json:"send_name"`
	SendID      int64     `json:"send_id"`
	Status      int       `json:"status"`
	SendPicture string    `json:"send_pictures"`
}

type ListAppResponse struct {
	Msg  string  `json:"msg" example:"success"`
	List []Apply `json:"list"`
}
type RefuseFriend struct {
	Account_id    int64  `json:"account_id" binding:"required"`
	Account_email string `json:"-"`
	UserID        int64  `json:"-"`
	AppliID       int64  `json:"-"`
}

type AcceptFriend struct {
	AppliID    int64 `json:"-"`
	Account_id int64 `json:"account_id" binding:"required"`
	UserID     int64 `json:"-"`
}

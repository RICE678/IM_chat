package models

type AppliSearch struct {
	UserID    int64  `json:"-"`
	SendEmail string `json:"send_email" binding:"required" form:"send_email"`
	SendID    int64  `json:"-"`
	Msg       string `json:"msg" binding:"required" form:"msg"`
}

package models

import "time"

type AppliSearch struct {
	UserID    int64  `json:"-"`
	SendEmail string `json:"send_email" binding:"required" form:"send_email"`
	SendID    int64  `json:"-"`
	Msg       string `json:"msg" binding:"required" form:"msg"`
}

type Apply struct {
	SendEmail string    `json:"send_email"`
	Msg       string    `json:"msg"`
	Time      time.Time `json:"time"`
	SendName  string    `json:"send_name"`
	SendID    int64     `json:"send_id"`
	Status    int       `json:"status"`
}

type Applies *[]Apply

// ListAppResponse 我发出的申请列表接口响应
type ListAppResponse struct {
	Msg  string  `json:"msg" example:"success"`
	List []Apply `json:"list"`
}

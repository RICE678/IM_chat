package models

type RegisterServer struct {
	Email      string `json:"email" form:"email" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required"`
	RePassword string `json:"re_password" form:"re_password" binding:"required"`
	Code       string `json:"code" form:"code" binding:"required"`
}

type LoginServer struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	UserID   int64  `json:"-"`
	Token    string `json:"-"`
}

type LoginResponse struct {
	Msg    string `json:"msg"`
	Token  string `json:"token,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
}

type ReUpdate struct {
	Code        string `json:"code" form:"code" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
	Email       string `json:"-"`
	UserID      int64  `json:"-"`
}

type ReEmail struct {
	NewEmail string `json:"new_email" form:"new_email" binding:"required"`
	Code     string `json:"code" form:"code" binding:"required"`
	Email    string `json:"-"`
	UserID   int64  `json:"-"`
	Password string `json:"-"`
}

type UserMain struct {
	Name      string `form:"name" json:"name"`
	Gender    int    `form:"gender" json:"gender"` //0为男 1为女 2为未知
	Signature string `form:"signature" json:"signature"`
	UserID    int64  `json:"-"`
	Picture   string `form:"picture" json:"picture"`
	Email     string `form:"email" json:"email"`
}

type UserMainReturn struct {
	User UserMain `json:"user,omitempty"`
	Err  string   `json:"err"`
}

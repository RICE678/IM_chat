package userserve

import (
	"IM_chat/dao"
	redis2 "IM_chat/dao/redisdao"
	"IM_chat/dao/sql"
	email2 "IM_chat/logic/email"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/jwt"
	"IM_chat/pkg/redis"
	"context"
	"fmt"
	"time"
)

// RegisterDetail 注册详情
func RegisterDetail(user *models.RegisterServer) string {
	if sql.IsRegisterEmail(user.Email) {
		return errcode.Msg(errcode.HasRegister)
	}
	if user.Password != user.RePassword {
		return errcode.Msg(errcode.PasswordNotMatch)
	}
	code := redis.RDB.Get(context.Background(), "email:"+user.Email).Val()
	if code != user.Code {
		return errcode.Msg(errcode.CodeError)
	}
	user.Password = dao.Md5(user.Password)
	if err := sql.AddRegister(user.Email, user.Password); err != nil {
		return errcode.Msg(errcode.NotRegister)
	}
	redis.RDB.Del(context.Background(), "email:"+user.Email)
	redis.RDB.Set(context.Background(), "login:"+user.Email, user.Password, time.Hour*24)
	return errcode.Msg(errcode.SUCCESS)
}

func LoginDetail(user *models.LoginServer) (string, string, int64) {
	user.Password = dao.Md5(user.Password)
	exists, err := redis.RDB.Exists(context.Background(), "login:"+user.Email).Result()
	if err != nil {
		return errcode.Msg(errcode.ERROR), "", 0
	}
	if exists == 0 {
		if !sql.IsRegisterEmail(user.Email) {
			return errcode.Msg(errcode.NullEmail), "", 0
		}
		if err = sql.Login(user.Email, user.Password); err != nil {
			return errcode.Msg(errcode.PasswordError), "", 0
		}
	} else {
		password := redis.RDB.Get(context.Background(), "login:"+user.Email).Val()
		if password != user.Password {
			return errcode.Msg(errcode.PasswordError), "", 0
		}
	}
	user.UserID, err = sql.SearchID(user.Email)
	if err != nil {
		return errcode.Msg(errcode.ERROR), "", 0
	}
	if err = redis2.LoginRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR), "", 0
	}
	token, err := jwt.SetToken(user.UserID)
	if err != nil {
		return errcode.Msg(errcode.ERROR), "", 0
	}
	user.Token = token
	return errcode.Msg(errcode.SUCCESS), token, user.UserID
}

func UpdateDetail(user *models.ReUpdate) string {
	var err error
	if user.Email, err = sql.SearchEmail(user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	codes := redis.RDB.Get(context.Background(), "email:"+user.Email).Val()
	if codes != user.Code {
		return errcode.Msg(errcode.CodeError)
	}
	user.NewPassword = dao.Md5(user.NewPassword)
	redis.RDB.Del(context.Background(), "email:"+user.Email)
	if err = redis2.ReLoginRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = sql.ReSetPassword(user.UserID, user.NewPassword); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

func ReEmail(user *models.ReEmail) string {
	var err error
	if user.Email, err = sql.SearchEmail(user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	codes := redis.RDB.Get(context.Background(), "email:"+user.Email).Val()
	if codes != user.Code {
		return errcode.Msg(errcode.CodeError)
	}
	if user.Password, err = sql.SearchPassword(user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = sql.ReSetEmail(user.UserID, user.NewEmail); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	redis.RDB.Del(context.Background(), "email:"+user.Email)
	if err = redis2.ReEmailRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

func CreateAccountDetails(user *models.UserMain2) string {
	if err := sql.UpdateUserMain(user.UserID, user.Name, user.Gender, user.Signature, user.PictureID); err != nil {
		return fmt.Sprintf("服务器出错: %v", err)
	}
	return errcode.Msg(errcode.SUCCESS)
}
func SearchAccountDetails(userID int64) (u models.UserMain, err string) {
	user, err2 := sql.SearchUserMain(userID)
	if err2 != nil {
		err = errcode.Msg(errcode.ERROR)
		return
	}
	u.Picture, err2 = sql.SearchPicture(userID)
	if err2 != nil {
		err = errcode.Msg(errcode.ERROR)
		return
	}
	u.PictureID, _ = sql.SearchPictureID(userID)
	u.Email = user.Email
	u.UserID = userID
	u.Name = user.Name
	u.Gender = user.Gender
	u.Signature = user.Signature
	err = errcode.Msg(errcode.SUCCESS)
	return
}

func DelUserDetail(userID int64) string {
	if err := redis2.LogoutRedis(userID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err := sql.DeleteUser(userID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

func ReCodeSend(userID int64) string {
	var email string
	var err error
	if email, err = sql.SearchEmail(userID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if !dao.VerifyEmailFormat(email) {
		return errcode.Msg(errcode.InvalidEmail)
	}
	if redis.RDB.Get(context.Background(), "send-email:"+email).Val() != "" {
		return errcode.Msg(errcode.HasSendCode)
	}
	code := dao.GetConfirmCode()
	if err = email2.SendConfirmMessage(email, code); err != nil {
		return errcode.Msg(errcode.DontSendCode)
	}
	redis.RDB.Set(context.Background(), "email:"+email, code, time.Minute*30)
	redis.RDB.Set(context.Background(), "send-email:"+email, code, time.Minute*1)
	return errcode.Msg(errcode.SUCCESS)
}

func SearchPictures() (ps []models.Pictures, err string) {
	var err2 error
	ps, err2 = sql.ShowPicture()
	if err2 != nil {
		err = errcode.Msg(errcode.ERROR)
		return
	}
	return
}

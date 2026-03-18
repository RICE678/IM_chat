package userserve

import (
	"IM_chat/dao"
	redis2 "IM_chat/dao/redisdao"
	"IM_chat/dao/sql"
	"IM_chat/initialize/redis"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/jwt"
	"context"
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

func LoginDetail(user *models.LoginServer) string {
	user.Password = dao.Md5(user.Password)
	exists, err := redis.RDB.Exists(context.Background(), "login:"+user.Email).Result()
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if exists == 0 {
		if !sql.IsRegisterEmail(user.Email) {
			return errcode.Msg(errcode.NullEmail)
		}
		if err = sql.Login(user.Email, user.Password); err != nil {
			return errcode.Msg(errcode.PasswordError)
		}
	} else {
		password := redis.RDB.Get(context.Background(), "login:"+user.Email).Val()
		if password != user.Password {
			return errcode.Msg(errcode.PasswordError)
		}
	}
	user.UserID, err = sql.SearchID(user.Email)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redis2.LoginRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	token, err := jwt.SetToken(user.UserID)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	user.Token = token
	return errcode.Msg(errcode.SUCCESS)
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
	var password string
	if password, err = sql.SearchEmail(user.UserID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	user.Password = dao.Md5(password)
	if err := redis2.ReEmailRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = sql.ReSetEmail(user.UserID, user.Email); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redis2.ReEmailRedis(user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

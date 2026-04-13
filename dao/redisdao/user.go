package redisdao

import (
	"IM_chat/models"
	"IM_chat/pkg/redis"
	"context"
	"strconv"
	"time"
)

func LoginRedis(user *models.LoginServer) error {
	if err := redis.RDB.Del(context.Background(), "login:"+user.Email).Err(); err != nil {
		return err
	}
	if err := redis.RDB.Set(context.Background(), "login:"+user.Email, user.Password, time.Hour*48).Err(); err != nil {
		return err
	}
	if err := redis.RDB.Set(context.Background(), "online:"+strconv.FormatInt(user.UserID, 10), 1, time.Second*30).Err(); err != nil {
		return err //登录时设置30s过期 此后每15秒心跳检测是否在线
	}
	return nil
}

func ReLoginRedis(user *models.ReUpdate) error {
	if err := redis.RDB.Del(context.Background(), "login:"+user.Email).Err(); err != nil {
		return err
	}
	if err := redis.RDB.Set(context.Background(), "login:"+user.Email, user.NewPassword, time.Hour*48).Err(); err != nil {
		return err
	}
	if err := redis.RDB.Set(context.Background(), "online:"+strconv.FormatInt(user.UserID, 10), 1, time.Second*30).Err(); err != nil {
		return err //登录时设置30s过期 此后每15秒心跳检测是否在线
	}
	return nil
}
func ReEmailRedis(user *models.ReEmail) error {
	if err := redis.RDB.Del(context.Background(), "login:"+user.Email).Err(); err != nil {
		return err
	}
	if err := redis.RDB.Set(context.Background(), "login:"+user.NewEmail, user.Password, time.Hour*48).Err(); err != nil {
		return err
	}
	return nil
}

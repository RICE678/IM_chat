package redisdao

import (
	"IM_chat/dao/sql"
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/redis"
	"context"
	"errors"
	"strconv"
	"time"
)

func RefreshOnline(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	onlineKey := "online:" + strconv.FormatInt(userID, 10)
	exists, err := redis.RDB.Exists(ctx, onlineKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New(errcode.Msg(errcode.UserNotLogin))
	}
	return redis.RDB.Expire(ctx, onlineKey, 30*time.Second).Err()
}

func IsUserOnline(userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	onlineKey := "online:" + strconv.FormatInt(userID, 10)
	exists, err := redis.RDB.Exists(ctx, onlineKey).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func LogoutRedis(userID int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := strconv.FormatInt(userID, 10)
	onlineKey := "online:" + id
	email, err := sql.SearchEmail(userID)
	if err != nil {
		return err
	}
	loginKey := "login:" + email
	_, err = redis.RDB.Del(ctx, loginKey, onlineKey).Result()
	if err != nil {
		return err
	}
	return nil
}

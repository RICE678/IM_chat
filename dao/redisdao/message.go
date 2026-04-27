package redisdao

import (
	"IM_chat/pkg/errcode"
	redisinit "IM_chat/pkg/redis"
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const (
	unread    = "message:unread:%d:%d"
	idlockkey = "read:lock:%d:%d"
)

func unreadKey(userID, friendID int64) string {
	return fmt.Sprintf(unread, userID, friendID)
}

func InitUnreadCount(userID, friendID int64) error {
	ctx := context.Background()
	return redisinit.RDB.Set(ctx, unreadKey(userID, friendID), 0, applyInboxTTL).Err()
}

func IncrUnreadCount(userID, friendID int64) error {
	ctx := context.Background()
	key := unreadKey(userID, friendID)
	if err := redisinit.RDB.Incr(ctx, key).Err(); err != nil {
		return err
	}
	if err := redisinit.RDB.Expire(ctx, key, applyInboxTTL).Err(); err != nil {
		return err
	}
	return nil
}

func GetUnreadCount(userID, friendID int64) (int, error) {
	ctx := context.Background()
	key := unreadKey(userID, friendID)
	num, err := redisinit.RDB.Get(ctx, key).Int()
	return num, err
}

func DelUnreadCount(userID, friendID int64) error {
	ctx := context.Background()
	key := unreadKey(userID, friendID)
	err := redisinit.RDB.Del(ctx, key).Err()
	return err
}

func CheckKeyMessage(userID, friendID int64) string {
	idkey := fmt.Sprintf(idlockkey, userID, friendID)
	ok, err := redisinit.RDB.SetNX(context.Background(), idkey, "1", 5*time.Second).Result()
	if err != nil {
		zap.L().Error("redis setnx error", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	if !ok {
		zap.L().Info("duplicate read request,skip", zap.Int64("userId", userID), zap.Int64("friendId", friendID))
		return errcode.Msg(errcode.SUCCESS)
	}
	return ""
}

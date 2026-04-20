package redisdao

import (
	redisinit "IM_chat/pkg/redis"
	"context"
	"fmt"
)

const (
	unread = "message:unread:%d:%d"
)

func CreateUnRead(userID, friendID int64, num int) error {
	ctx := context.Background()
	key := fmt.Sprintf(unread, userID, friendID)
	if err := redisinit.RDB.Set(ctx, key, num, applyInboxTTL).Err(); err != nil {
		return err
	}
	return nil
}

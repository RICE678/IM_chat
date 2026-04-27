package redisdao

import (
	redisinit "IM_chat/pkg/redis"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	applyKeyLock  = "apply:lock:%d:%d"
	applyKeySend  = "apply:send:%d"
	applyKeyInbox = "apply:inbox:%d"
	applyKeyInfo  = "apply:info:%d"
	applyPair     = "apply:pair:%d:%d"

	// 防重复锁：10 分钟内不允许同一人对同一人重复申请
	applyLock = 10 * time.Minute

	// 收件箱/发件箱 ZSET 默认过期：7 天
	applyInboxTTL = 7 * 24 * time.Hour
)

// SetApplyLock 防重复：同一人对同一人短期内不能重复申请
func SetApplyLock(fromID, toID int64) (bool, error) {
	ctx := context.Background()
	lockKey := fmt.Sprintf(applyKeyLock, fromID, toID)
	res, err := redisinit.RDB.SetNX(ctx, lockKey, "1", applyLock).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return res, nil
}

// PushApplyInbox 申请入库后，写入收件箱索引
func PushApplyInbox(toID, applyID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeyInbox, toID)
	score := float64(time.Now().UnixMilli())
	if err := redisinit.RDB.ZAdd(ctx, key, redis.Z{Score: score, Member: applyID}).Err(); err != nil {
		return err
	}
	return redisinit.RDB.Expire(ctx, key, applyInboxTTL).Err()
}

// PushApplySent 申请入库后，写入发件箱索引
func PushApplySent(fromID, applyID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeySend, fromID)
	score := float64(time.Now().UnixMilli())
	if err := redisinit.RDB.ZAdd(ctx, key, redis.Z{Score: score, Member: applyID}).Err(); err != nil {
		return err
	}
	return redisinit.RDB.Expire(ctx, key, applyInboxTTL).Err()
}

// SetApplyInfo 保存申请详情到Redis
func SetApplyInfo(applyID, fromID, toID int64, msg string) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeyInfo, applyID)
	now := time.Now().UnixMilli()
	if err := redisinit.RDB.HSet(ctx, key,
		"apply_id", applyID,
		"from_id", fromID,
		"to_id", toID,
		"msg", msg,
		"status", 0,
		"create_time_ms", now,
	).Err(); err != nil {
		return err
	}
	return redisinit.RDB.Expire(ctx, key, applyInboxTTL).Err()
}

// SetApplyCache 一次性写入申请相关缓存
func SetApplyCache(applyID, fromID, toID int64, msg string) error {
	if err := PushApplyInbox(toID, applyID); err != nil {
		return err
	}
	if err := PushApplySent(fromID, applyID); err != nil {
		return err
	}
	if err := SetApplyInfo(applyID, fromID, toID, msg); err != nil {
		return err
	}
	if err := SetApplyPair(applyID, fromID, toID); err != nil {
		return err
	}
	return nil
}

// SetApplyPair 建立接收人-发送人到申请ID的映射。
func SetApplyPair(applyID, senderID, receiverID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyPair, receiverID, senderID)
	if err := redisinit.RDB.Set(ctx, key, applyID, applyInboxTTL).Err(); err != nil {
		return err
	}
	return nil
}

// GetApplyPair 根据 (receiverID, senderID) 获取申请ID。
func GetApplyPair(receiverID, senderID int64) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf(applyPair, receiverID, senderID)
	val, err := redisinit.RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DelApplyPair(receiverID, senderID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyPair, receiverID, senderID)
	return redisinit.RDB.Del(ctx, key).Err()
}

func SetApplyStatus(applyID int64, status int) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeyInfo, applyID)
	if err := redisinit.RDB.HSet(ctx, key,
		"status", status,
	).Err(); err != nil {
		return err
	}
	return redisinit.RDB.Expire(ctx, key, applyInboxTTL).Err()
}

func GetApplySend(fromID int64) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeySend, fromID)
	res, err := redisinit.RDB.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetApplyInboxIDs 获取某人收件箱中的申请 ID 列表
func GetApplyInboxIDs(toID int64, limit int64) ([]int64, error) {
	if limit <= 0 {
		limit = 20
	}
	ctx := context.Background()
	key := fmt.Sprintf(applyKeyInbox, toID)
	ids, err := redisinit.RDB.ZRevRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}
	result := make([]int64, 0, len(ids))
	for _, s := range ids {
		id, _ := strconv.ParseInt(s, 10, 64)
		if id > 0 {
			result = append(result, id)
		}
	}
	return result, nil
}

func DelApplyFromInbox(toID, applyID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf(applyKeyInbox, toID)
	return redisinit.RDB.ZRem(ctx, key, applyID).Err()
}

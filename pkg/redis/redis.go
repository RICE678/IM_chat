package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis() (err error) {
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	if host == "" || port == 0 {
		return fmt.Errorf("redis config missing: redis.host / redis.port in config.yaml")
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		return err
	}
	zap.L().Info("redis connected", zap.String("addr", fmt.Sprintf("%s:%d", host, port)), zap.Int("db", viper.GetInt("redis.db")))
	return nil
}

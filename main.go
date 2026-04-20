package main

import (
	_ "IM_chat/docs"
	"IM_chat/logic/mq"
	"IM_chat/pkg/errcode"
	kafka "IM_chat/pkg/kafkapkg"
	"IM_chat/pkg/snowflake"

	//"IM_chat/logic/mq"
	"IM_chat/middlewares"
	"IM_chat/pkg/config"
	//"IM_chat/pkg/errcode"
	//"IM_chat/pkg/kafkapkg"
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/redis"
	"IM_chat/routes"
	"IM_chat/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title IM Chat API
// @version 2.0
// @description IM Chat backend API documentation.
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func findConfigPath() string {
	possiblePaths := []string{
		filepath.Join("config", "config.yaml"),
		filepath.Join("..", "config", "config.yaml"),
	}
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return filepath.Join("config", "config.yaml")
}

func main() {
	cfgPath := findConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Fail to load config file: %v", err)
	}
	if err := settings.InitSettings(cfgPath); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	if err := middlewares.InitLogger(); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	if err := mysql.InitMysql(); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	if err := redis.InitRedis(); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	kafka.InitKafka(&cfg.Kafka)
	if initMsg := kafka.InitProducer(&cfg.Kafka); initMsg != errcode.Msg(errcode.SUCCESS) {
		fmt.Printf("init kafka failed, msg:%s\n", initMsg)
	}
	groupID := int64(cfg.Kafka.GroupID)
	consumer, errStr := kafka.NewPersistConsumer(
		cfg.Kafka.Brokers,
		groupID,
		[]string{kafka.TopicPrivateMsg},
		mq.HandleKafkaMessage,
	)
	if errStr != errcode.Msg(errcode.SUCCESS) {
		fmt.Printf("init kafka failed, msg:%s\n", errStr)
	} else {
		go consumer.Start()
		defer consumer.Stop()
	}
	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("init snowflake failed,err:%v\n", err)
		return
	}

	r := routes.Setup()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exit")
}

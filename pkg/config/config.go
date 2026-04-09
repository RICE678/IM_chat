package config

import (
	"IM_chat/models"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Kafka models.KafkaConfig `json:"kafka"`
}

var AppConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file:%w", err)
	}
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file:%w", err)
	}
	AppConfig = config
	log.Printf("Config loaded successfully from %s", configPath)
	return config, nil
}

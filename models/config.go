package models

var Conf = new(Config)

type Config struct {
	*AppConfig `yaml:"app"`
	Kafka      KafkaConfig `yaml:"kafka"`
}
type AppConfig struct {
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	StartTime string `json:"start_time"`
	MachineID int64  `json:"machine_id"`
	Port      int    `json:"port"`
}

type KafkaConfig struct {
	Brokers []string    `json:"brokers"`
	Topic   string      `json:"topic"`
	Topics  KafkaTopics `json:"topics"`
	GroupID int         `json:"group_id"`
}
type KafkaTopics struct {
	GroupMsgRaw     string `json:"group_msg_raw"`
	GroupMsgPublish string `json:"group_msg_publish"`
	PrivateMsg      string `json:"private_msg"`
	ReadEvent       string `json:"read_event"`
}

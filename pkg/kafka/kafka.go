package kafka

import "IM_chat/models"

func InitKafka(cfg *models.KafkaConfig) {
	if cfg.Topics.GroupMsgRaw != "" {
		TopicGroupMsgRaw = "group_msg_raw"
	}
	if cfg.Topics.GroupMsgPublish != "" {
		TopicGroupMsgPublish = "group_msg_publish"
	}
	if cfg.Topics.PrivateMsg != "" {
		TopicPrivateMsg = cfg.Topics.PrivateMsg
	}
	if cfg.Topics.ReadEvent != "" {
		TopicReadEvent = cfg.Topics.ReadEvent
	}
}

var (
	TopicGroupMsgRaw     = "group-msg-raw"
	TopicGroupMsgPublish = "group-msg-publish"
	TopicPrivateMsg      = "private-msg"
	TopicReadEvent       = "read-event"
)

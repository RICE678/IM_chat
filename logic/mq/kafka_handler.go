package mq

import (
	"IM_chat/logic/ws"
	"IM_chat/models"
	"encoding/json"
	"go.uber.org/zap"
)

func HandleKafkaMessage(topic string, key []byte, value []byte) error {
	var msg models.WsMsg
	if err := json.Unmarshal(value, &msg); err != nil {
		zap.L().Error("unmarshal msg failed", zap.Error(err))
		return err
	}
	ok := ws.GlobalManager.Send(msg.ReceiverID, &msg)
	zap.L().Info("kafka consumed",
		zap.String("topic", topic),
		zap.Int64("receiver", msg.ReceiverID),
		zap.Bool("delivered", ok))
	return nil
}

package ws

import (
	"IM_chat/dao/redisdao"
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	kafka2 "IM_chat/pkg/kafkapkg"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func HandleMessage(sender *Client, msg *models.WsMsg) {
	switch msg.Type {
	case "private":
		HandlePrivateMessage(sender, msg)
	case "text":
		HandleText(sender, msg)
	default:
		zap.L().Error("unknown type", zap.String("type", msg.Type))
	}
}

func HandleText(sender *Client, msg *models.WsMsg) {
	if msg.ReceiverID == 0 || msg.Msg == "" {
		return
	}
	msg.Timestamp = time.Now().UnixMilli()
	key := strconv.FormatInt(msg.ReceiverID, 10)
	errStr := kafka2.Publish(kafka2.TopicPrivateMsg, key, msg)
	if errStr != errcode.Msg(errcode.SUCCESS) {
		zap.L().Error("publish msg failed", zap.String("err", errStr))
	}
	ack := &models.WsMsg{
		Type:       "ack",
		Timestamp:  msg.Timestamp,
		Msg:        "ok",
		ReceiverID: msg.ReceiverID,
		MsgType:    msg.MsgType,
	}
	sender.Send(ack)
}

func HandlePrivateMessage(sender *Client, msg *models.WsMsg) {
	zap.L().Info("private msg in",
		zap.Int64("from", msg.SenderID),
		zap.Int64("to", msg.ReceiverID),
		zap.String("content", msg.Msg))
	if msg.ReceiverID == 0 || msg.Msg == "" {
		zap.L().Warn("private msg dropped: empty receiver or content")
		return
	}
	msg.Timestamp = time.Now().UnixMilli()
	key := strconv.FormatInt(msg.ReceiverID, 10)
	errStr := kafka2.Publish(kafka2.TopicPrivateMsg, key, msg)
	zap.L().Info("kafka publish result", zap.String("err", errStr))
	if errStr != errcode.Msg(errcode.SUCCESS) {
		zap.L().Error("publish msg failed", zap.String("err", errStr))
	}
	ack := &models.WsMsg{
		Type:       "ack",
		Timestamp:  msg.Timestamp,
		Msg:        "ok",
		ReceiverID: msg.ReceiverID,
		MsgType:    msg.MsgType,
	}
	zap.L().Info("sending ack to sender", zap.Int64("sender", sender.UserID))
	sender.Send(ack)
}

func HistoryMain(user *models.HistoryMsg) string {
	var err error
	if user.Msg, err = sql.GetHistory(user.SenderID, user.ReceiverID, user.Page, user.Size); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

func ReadMain(user *models.ReadContact) string {
	if user.FriendID <= 0 {
		return errcode.Msg(errcode.InvalidParams)
	}
	if err1 := redisdao.CheckKeyMessage(user.FriendID, user.UserID); err1 != "" {
		return err1
	}
	count, err := redisdao.GetUnreadCount(user.FriendID, user.UserID)
	if err != nil && count == 0 {
		return errcode.Msg(errcode.SUCCESS)
	}
	err = sql.ReadContact(user.UserID, user.FriendID)
	if err != nil {
		zap.L().Error("change read contact failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.DelUnreadCount(user.FriendID, user.UserID); err != nil {
		zap.L().Error("delete unread count failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.InitUnreadCount(user.FriendID, user.UserID); err != nil {
		zap.L().Error("init unread count failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	err = sql.ChangeRead(user.FriendID, user.UserID)
	if err != nil {
		zap.L().Error("change read message failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}

	return errcode.Msg(errcode.SUCCESS)
}

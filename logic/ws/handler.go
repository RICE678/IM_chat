package ws

import (
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

func HistoryMain(user *models.HistoryMsg) string {
	var err error
	if user.Msg, err = sql.GetHistory(user.SenderID, user.ReceiverID, user.Page, user.Size); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

//	func UnReadMain(userID int64) ([]models.MainFriend, string) {
//		rows, err := sql.SearchUnRead(userID)
//		if err != nil {
//			return nil, errcode.Msg(errcode.ErrorForList)
//		}
//		list := make([]models.MainFriend, 0, len(rows)+16)
//		for _, r := range rows {
//			list = append(list, models.MainFriend{
//				FriendID:    r.FriendID,
//				LastMsgTime: r.Last_msg_time,
//				Unread:      r.Unread_contact,
//			})
//		}
//		return list, errcode.Msg(errcode.SUCCESS)
//	}
func ReadMain(user *models.ReadContact) string {
	if user.FriendID <= 0 {
		return errcode.Msg(errcode.InvalidParams)
	}
	err := sql.ReadContact(user.UserID, user.FriendID)
	if err != nil {
		zap.L().Error("change read contact failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	err = sql.ChangeRead(user.FriendID, user.UserID)
	if err != nil {
		zap.L().Error("change read message failed", zap.Error(err))
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

package sql

import (
	"IM_chat/models"
	"IM_chat/pkg/mysql"
	stdsql "database/sql"
	"time"
)

type messageRow struct {
	ID         int64           `db:"id"`
	SenderID   int64           `db:"sender_id"`
	ReceiverID int64           `db:"receiver_id"`
	MsgType    int             `db:"msg_type"`
	Content    string          `db:"content"`
	SendTime   stdsql.NullTime `db:"send_time"`
}

func SaveMessage(msg *models.ChatMsg) error {
	if msg.CreateTime.IsZero() {
		msg.CreateTime = time.Now()
	}
	if msg.Context == "" {
		msg.Context = msg.Msg
	}
	if msg.Msg == "" {
		msg.Msg = msg.Context
	}
	_, err := mysql.MessagesDB().Exec("insert into messages(id,sender_id,receiver_id,send_status,msg_type,content,is_delete,read_status,send_time)values(?,?,?,?,?,?,?,?,?)",
		msg.ID, msg.UserID, msg.ReceiverID, 1, msg.MsgType, msg.Context, 0, 0, msg.CreateTime)
	return err
}

func toChatMsg(r messageRow) *models.ChatMsg {
	msg := &models.ChatMsg{
		ID:         r.ID,
		UserID:     r.SenderID,
		ReceiverID: r.ReceiverID,
		MsgType:    r.MsgType,
		Msg:        r.Content,
		Context:    r.Content,
	}
	if r.SendTime.Valid {
		msg.CreateTime = r.SendTime.Time
		msg.Timestamp = r.SendTime.Time.UnixMilli()
	}
	return msg
}

func GetHistory(userA, userB int64, page, size int) (msgs []*models.ChatMsg, err error) {
	offset := (page - 1) * size
	var rows []messageRow
	err = mysql.MessagesDB().Select(
		&rows,
		"select id,sender_id,receiver_id,msg_type,content,send_time from messages where ((sender_id=? and receiver_id=?) or (sender_id=? and receiver_id=?)) and send_status=1 order by id desc limit ? offset ?",
		userA, userB, userB, userA, size, offset,
	)
	if err != nil {
		return nil, err
	}
	msgs = make([]*models.ChatMsg, 0, len(rows))
	for _, r := range rows {
		msgs = append(msgs, toChatMsg(r))
	}
	return msgs, nil
}

func GetUnreadMessages(receiverID int64) (msgs []*models.ChatMsg, err error) {
	var rows []messageRow
	err = mysql.MessagesDB().Select(
		&rows,
		"select id,sender_id,receiver_id,msg_type,content,send_time from (select id,sender_id,receiver_id,msg_type,content,send_time from messages where (receiver_id=? or sender_id=?) and send_status=1 and is_delete=0 order by id desc limit 500) t order by id asc",
		receiverID, receiverID)
	msgs = make([]*models.ChatMsg, 0, len(rows))
	for _, r := range rows {
		msgs = append(msgs, toChatMsg(r))
	}
	return msgs, nil
}

func InsertUnRead(msg *models.WsMsg) error {
	_, err := mysql.DB().Exec("update contact set unread_contact=unread_contact+1 where user_id=? and friend_id=?", msg.ReceiverID, msg.SenderID)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = mysql.DB().Exec("update contact set last_msg_time=? where user_id=? and friend_id=?", now, msg.ReceiverID, msg.SenderID)
	if err != nil {
		return err
	}
	_, err = mysql.DB().Exec("update contact set last_msg_time=? where user_id=? and friend_id=?", now, msg.SenderID, msg.ReceiverID)
	return err
}

func ChangeRead(userA, userB int64) error {
	_, err := mysql.MessagesDB().Exec("update messages set read_status=1 where send_status=1 and is_delete=0 and sender_id=? and receiver_id=? and read_status=0", userA, userB)
	return err
}

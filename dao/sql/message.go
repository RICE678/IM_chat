package sql

import (
	"IM_chat/models"
	"IM_chat/pkg/mysql"
	"time"
)

type messageRow struct {
	ID         int64  `db:"id"`
	SenderID   int64  `db:"sender_id"`
	ReceiverID int64  `db:"receiver_id"`
	MsgType    int    `db:"msg_type"`
	Content    string `db:"content"`
}

func SaveMessage(msg *models.ChatMsg) error {
	_, err := mysql.DB().Exec("insert into messages(id,sender_id,receiver_id,send_status,msg_type,content,is_delete,read_status)values(?,?,?,?,?,?,?,?)",
		msg.ID, msg.UserID, msg.ReceiverID, 1, msg.MsgType, msg.Context, 0, 0)
	return err
}

func GetHistory(userA, userB int64, page, size int) (msgs []*models.ChatMsg, err error) {
	offset := (page - 1) * size
	var rows []messageRow
	err = mysql.DB().Select(
		&rows,
		"select id,sender_id,receiver_id,msg_type,content from messages where ((sender_id=? and receiver_id=?) or (sender_id=? and receiver_id=?)) and send_status=1 order by id desc limit ? offset ?",
		userA, userB, userB, userA, size, offset,
	)
	if err != nil {
		return nil, err
	}
	msgs = make([]*models.ChatMsg, 0, len(rows))
	for _, r := range rows {
		msgs = append(msgs, &models.ChatMsg{
			ID:         r.ID,
			UserID:     r.SenderID,
			ReceiverID: r.ReceiverID,
			MsgType:    r.MsgType,
			Context:    r.Content,
		})
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
	_, err := mysql.DB().Exec("update messages set read_status=1 where send_status=1 and is_delete=0 and sender_id=? and receiver_id=? and read_status=0", userA, userB)
	return err
}

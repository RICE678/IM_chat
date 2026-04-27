package sql

import (
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/snowflake"
	"fmt"
	"time"
)

type ContactRow struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	FriendID       int64     `db:"friend_id"`
	Add_time       time.Time `db:"add_time"`
	Friend_status  int       `db:"friend_status"`
	Is_del         int       `db:"is_del"`
	Last_msg_time  time.Time `db:"last_msg_time"`
	Unread_contact int       `db:"unread_contact"`
}

func ensureContact(userID, friendID int64) error {
	if IsFriend(userID, friendID) {
		return nil
	}
	id := snowflake.Generate()
	_, err := mysql.DB().Exec("insert into contact(id,user_id,friend_id,friend_status)values(?,?,?,?)", id, userID, friendID, 1)
	return err
}

func CreateContact(userID, friendID int64) error {
	if userID <= 0 || friendID <= 0 {
		return fmt.Errorf("invalid contact ids: userID=%d friendID=%d", userID, friendID)
	}
	if err := ensureContact(userID, friendID); err != nil {
		return err
	}
	return ensureContact(friendID, userID)
}

func SearchUnRead(userID int64) ([]ContactRow, error) {
	var row1, row2 []ContactRow
	err := mysql.DB().Select(&row1, "select * from contact where user_id=? and friend_status=1 and is_del=0 and unread_contact>0 order by last_msg_time DESC ", userID)
	if err != nil {
		return nil, err
	}
	err = mysql.DB().Select(&row2, "select * from contact where user_id=? and friend_status=1 and is_del=0 and unread_contact=0 order by last_msg_time DESC ", userID)
	if err != nil {
		return nil, err
	}
	rows := append(row1, row2...)
	return rows, nil
}

func ReadContact(userID, friendID int64) error {
	_, err := mysql.DB().Exec("update contact set unread_contact=0 where user_id=? and friend_id=?", userID, friendID)
	return err
}

func SearchContact(userID int64) ([]User, error) {
	var users []User
	err := mysql.DB().Select(&users, `
		select
			u.id,
			case
				when trim(coalesce(c.remark, '')) <> '' then trim(c.remark)
				else u.username
			end as username,
			COALESCE(c.remark, '') as remark,
			u.email,
			COALESCE(u.picture, ?) as picture
		from contact c
		join users u on u.id = c.friend_id
		where c.user_id = ?
			and c.friend_status = 1
		order by
			lower(
				case
					when trim(coalesce(c.remark, '')) <> '' then trim(c.remark)
					else u.username
				end
			) asc,
			u.id asc
	`, defaultPictureID, userID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func SearchShow(userID int64) ([]User, error) {
	var users []User
	err := mysql.DB().Select(&users, `
		select
			u.id,
			case
				when trim(coalesce(c.remark, '')) <> '' then trim(c.remark)
				else u.username
			end as username,
			COALESCE(c.remark, '') as remark,
			u.email,
			COALESCE(u.picture, ?) as picture
		from contact c
		join users u on u.id = c.friend_id
		where c.user_id = ?
			and c.friend_status = 1
			and c.is_del = 0
		order by
			case when c.unread_contact > 0 then 0 else 1 end asc,
			c.last_msg_time desc,
			u.id asc
	`, defaultPictureID, userID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteContact(userID, friendID int64) error {
	_, err := mysql.DB().Exec("update contact set is_del=1 and unread_contact=0 where user_id=? and friend_id=?", userID, friendID)
	return err
}

func IsFriend(userID, friendID int64) bool {
	var exists int
	err := mysql.DB().Get(
		&exists,
		"select exists(select 1 from contact where user_id=? and friend_id=? and friend_status=1)",
		userID, friendID,
	)
	if err != nil {
		return false
	}
	return exists == 1
}

func DelUnReadMessage(userID, friendID int64) error {
	_, err := mysql.MessagesDB().Exec(
		"update messages set read_status=1 where ((sender_id=? and receiver_id=?) or (sender_id=? and receiver_id=?)) and read_status=0",
		userID, friendID, friendID, userID,
	)
	return err
}

package sql

import (
	"IM_chat/models"
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

func CreateContact(user models.AcceptFriend) error {
	if user.UserID <= 0 || user.Account_id <= 0 {
		return fmt.Errorf("invalid contact ids: userID=%d account_id=%d", user.UserID, user.Account_id)
	}
	id := snowflake.Generate()
	_, err := mysql.DB().Exec("insert into contact(id,user_id,friend_id,friend_status)values(?,?,?,?)", id, user.UserID, user.Account_id, 1)
	if err != nil {
		return err
	}
	id1 := snowflake.Generate()
	_, err = mysql.DB().Exec("insert into contact(id,user_id,friend_id,friend_status)values(?,?,?,?)", id1, user.Account_id, user.UserID, 1)
	return err
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
			and c.is_del = 0
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

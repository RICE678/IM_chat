package sql

import (
	"IM_chat/models"
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/snowflake"
)

func CreateContact(user models.AcceptFriend) error {
	id := snowflake.Generate()
	_, err := mysql.DB().Exec("insert into contact(id,user_id,friend_id,friend_status)values(?,?,?,?)", id, user.UserID, user.Account_id, 1)
	return err
}

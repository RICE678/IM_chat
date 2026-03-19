package sql

import (
	"IM_chat/initialize/mysql"
	"IM_chat/pkg/snowflake"
	"time"
)

type User struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	Gender     int       `db:"gender"`
	Signature  string    `db:"signature"`
	CreateTime time.Time `db:"create_time"`
}

// IsRegisterEmail 是否注册过此邮箱
func IsRegisterEmail(email string) bool {
	var user User
	err := mysql.DB().Get(&user, "select id from users where email=?", email)
	if err != nil {
		return false
	}
	return user.ID > 0
}

func AddRegister(email, password string) error {
	_, err := mysql.DB().Exec("insert into users(id,password,email)values(?,?,?)", snowflake.Generate(), password, email)
	if err != nil {
		return err
	}
	return err
}

func Login(email, password string) error {
	var user User
	err := mysql.DB().Get(&user, "select id,password,email from users where email=? and password=?", email, password)
	if err != nil {
		return err
	}
	return nil
}

func SearchID(email string) (int64, error) {
	var user User
	err := mysql.DB().Get(&user, "select id from users where email=?", email)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func SearchEmail(id int64) (string, error) {
	var user User
	err := mysql.DB().Get(&user, "select email from users where id=?", id)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func SearchPassword(id int64) (string, error) {
	var user User
	err := mysql.DB().Get(&user, "select password from users where id=?", id)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

func ReSetEmail(id int64, email string) error {
	_, err := mysql.DB().Exec("update users set email=? where id=?", email, id)
	if err != nil {
		return err
	}
	return nil
}

func ReSetPassword(id int64, password string) error {
	_, err := mysql.DB().Exec("update users set password=? where id=?", password, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserMain(userid int64, name string, gender int, signature string) error {
	_, err := mysql.DB().Exec("update users set username=?, gender=?, signature=? where id=?", name, gender, signature, userid)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id int64) error {
	_, err := mysql.DB().Exec("update users set id_del=? where id=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

package sql

import (
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/snowflake"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID         int64     `db:"id"`
	Name       string    `db:"username"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	Gender     int       `db:"gender"`
	Signature  string    `db:"signature"`
	Picture    int       `db:"picture"`
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
	id := snowflake.Generate()
	idString := strconv.FormatInt(id, 10)
	suffix := idString
	if len(idString) > 6 {
		suffix = idString[len(idString)-6:]
	}
	name := "用户" + suffix
	_, err := mysql.DB().Exec("insert into users(id,username,password,email)values(?,?,?,?)", id, name, password, email)
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

func SearchName(id int64) (string, error) {
	var user User
	err := mysql.DB().Get(&user, "select username from users where id=?", id)
	if err != nil {
		return "", err
	}
	return user.Name, nil
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

func UpdateUserMain(userid int64, name string, gender int, signature string, picture string) error {
	var current User
	if err := mysql.DB().Get(&current, "select username, gender, signature, picture from users where id=?", userid); err != nil {
		return err
	}
	finalName := name
	if strings.TrimSpace(finalName) == "" {
		finalName = current.Name
	}
	finalSignature := signature
	if strings.TrimSpace(finalSignature) == "" {
		finalSignature = current.Signature
	}
	finalGender := gender
	if finalGender < 0 || finalGender > 2 {
		finalGender = current.Gender
	}
	finalPicture := current.Picture
	if strings.TrimSpace(picture) != "" {
		if err := mysql.DB().Get(&finalPicture, "select id from picture where web=?", picture); err != nil {
			return err
		}
	}
	_, err := mysql.DB().Exec(
		"update users set username=?, gender=?, signature=? ,picture=? where id=?",
		finalName, finalGender, finalSignature, finalPicture, userid,
	)
	if err != nil {
		return err
	}
	return nil
}

func SearchUserMain(userID int64) (u User, err error) {
	err = mysql.DB().Get(&u, "select username,email, gender, signature, picture from users where id=?", userID)
	return
}
func DeleteUser(id int64) error {
	_, err := mysql.DB().Exec("update users set id_del=? where id=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

func SearchPicture(id int64) (string, error) {
	var picid int
	var picture string
	err := mysql.DB().Get(&picid, "select picture from users where id=?", id)
	if err != nil {
		return "", err
	}
	err = mysql.DB().Get(&picture, "select web from picture where id=?", picid)
	if err != nil {
		return "", err
	}
	return picture, nil
}

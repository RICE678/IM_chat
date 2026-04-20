package sql

import (
	"IM_chat/models"
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/snowflake"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// defaultPictureID 系统默认头像（picture 表主键）；库里正常应始终存 1，COALESCE 仅兼容极少数 NULL 旧数据。
const defaultPictureID = 1

type User struct {
	ID         int64     `db:"id"`
	Name       string    `db:"username"`
	Remark     string    `db:"remark"`
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
	_, err := mysql.DB().Exec("insert into users(id,username,password,email,picture)values(?,?,?,?,?)", id, name, password, email, defaultPictureID)
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
	err := mysql.DB().Get(&user.ID, "select id from users where email=? and is_del=0", email)
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
func UpdateUserMain(userid int64, name string, gender *int, signature *string, pictureID *int) error {
	var current User
	if err := mysql.DB().Get(&current, "select username, gender, COALESCE(signature, '') as signature, COALESCE(picture, ?) as picture from users where id=?", defaultPictureID, userid); err != nil {
		return err
	}
	finalName := current.Name
	if strings.TrimSpace(name) != "" {
		finalName = strings.TrimSpace(name)
	}
	finalGender := current.Gender
	if gender != nil && *gender >= 0 && *gender <= 2 {
		finalGender = *gender
	}
	finalSignature := current.Signature
	if signature != nil {
		finalSignature = *signature
	}
	finalPicID := current.Picture
	if pictureID != nil {
		pid := *pictureID
		if pid < 1 {
			finalPicID = defaultPictureID
		} else {
			var found int
			if err := mysql.DB().Get(&found, "select id from picture where id=?", pid); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("invalid picture_id: %d", pid)
				}
				return err
			}
			finalPicID = found
		}
	}
	_, err := mysql.DB().Exec(
		"update users set username=?, gender=?, signature=? ,picture=? where id=?",
		finalName, finalGender, finalSignature, finalPicID, userid,
	)
	return err
}

func SearchUserMain(userID int64) (u User, err error) {
	err = mysql.DB().Get(&u, "select username,email, gender, COALESCE(signature, '') as signature, COALESCE(picture, ?) as picture from users where id=?", defaultPictureID, userID)
	return
}
func DeleteUser(id int64) error {
	_, err := mysql.DB().Exec("update users set is_del=? where id=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

func SearchPicture(id int64) (string, error) {
	var picid int
	err := mysql.DB().Get(&picid, "select COALESCE(picture, ?) from users where id=?", defaultPictureID, id)
	if err != nil {
		return "", err
	}
	var picture string
	err = mysql.DB().Get(&picture, "SELECT web FROM picture WHERE id = ?", picid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return picture, nil
}

func SearchPictureID(id int64) (int, error) {
	var picid int
	err := mysql.DB().Get(&picid, "select COALESCE(picture, ?) from users where id=?", defaultPictureID, id)
	if err != nil {
		return 0, err
	}
	return picid, nil
}

func ShowPicture() ([]models.Pictures, error) {
	var pictures []models.Pictures
	err := mysql.DB().Select(&pictures, "SELECT id, web FROM picture ORDER BY id ASC")
	return pictures, err
}

func SearchRemark(userID int64, friendID int64) (rename string, err error) {
	err = mysql.DB().Get(&rename, "select COALESCE(remark, '') from contact where user_id=? and friend_id=?", userID, friendID)
	return
}

func ChangeRemark(userID int64, friendID int64, remark string) (err error) {
	_, err = mysql.DB().Exec("update contact set remark=? where user_id=? and friend_id=?", remark, userID, friendID)
	return
}

package sql

import (
	"IM_chat/models"
	"IM_chat/pkg/mysql"
	"IM_chat/pkg/snowflake"
	"time"
)

// ApplyRow 申请记录（数据库映射）
type ApplyRow struct {
	ID         int64     `db:"id"`
	FromID     int64     `db:"from_id"`
	ToID       int64     `db:"to_id"`
	Remark     string    `db:"remark"`
	Status     int       `db:"status"`
	CreateTime time.Time `db:"create_time"`
}

func SetApply(user models.AppliSearch) (int64, error) {
	id := snowflake.Generate()
	_, err := mysql.DB().Exec("insert into apply(id,from_id,to_id,remark)values(?,?,?,?)", id, user.UserID, user.SendID, user.Msg)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ListApplySent 查询我发出的所有申请
func ListApplySent(fromID int64) ([]ApplyRow, error) {
	var rows []ApplyRow
	err := mysql.DB().Select(&rows, "select id,from_id,to_id,remark,status,create_time from apply where from_id=? order by create_time desc", fromID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func ListApplyTo(toID int64) ([]ApplyRow, error) {
	var rows []ApplyRow
	err := mysql.DB().Select(&rows, "select id,from_id,to_id,remark,status,create_time from apply where to_id=? order by create_time desc", toID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetApplyByID 按 ID 查单条申请
func GetApplyByID(applyID int64, fromID int64) (*ApplyRow, error) {
	var row ApplyRow
	err := mysql.DB().Get(&row, "select id,from_id,to_id,remark,status,create_time from apply where id=? and from_id=?", applyID, fromID)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func GetApplyByIDTo(applyID int64, toID int64) (*ApplyRow, error) {
	var row ApplyRow
	err := mysql.DB().Get(&row, "select id,from_id,to_id,remark,status,create_time from apply where id=? and to_id=?", applyID, toID)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func ChangeStatus(applyID int64, status int) error {
	_, err := mysql.DB().Exec("update apply set status=? where id=?", status, applyID)

	return err
}

// ChangeStatusByPair 按申请ID+收发双方更新状态，避免误更新或空更新
func ChangeStatusByPair(applyID, fromID, toID int64, status int) error {
	res, err := mysql.DB().Exec(
		"update apply set status=? where id=? and from_id=? and to_id=?",
		status, applyID, fromID, toID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return mysql.DB().Get(&ApplyRow{}, "select id,from_id,to_id,remark,status,create_time from apply where id=? and from_id=? and to_id=?", applyID, fromID, toID)
	}
	return nil
}

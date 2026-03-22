package sql

import (
	"IM_chat/initialize/mysql"
	"IM_chat/models"
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

// ListApplySent 查询我发出的所有申请（MySQL 主数据源，按 id 降序）
// 兼容仅有 id,from_id,to_id,remark 的表；若表有 status、create_time 可改用 select id,from_id,to_id,remark,status,create_time
func ListApplySent(fromID int64) ([]ApplyRow, error) {
	var rows []ApplyRow
	err := mysql.DB().Select(&rows, "select id,from_id,to_id,remark,status,create_time from apply where from_id=? order by create_time desc", fromID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetApplyByID 按 ID 查单条申请（用于 Redis 补全）
func GetApplyByID(applyID int64, fromID int64) (*ApplyRow, error) {
	var row ApplyRow
	err := mysql.DB().Get(&row, "select id,from_id,to_id,remark,0 as status,from_unixtime(0) as create_time from apply where id=? and from_id=?", applyID, fromID)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

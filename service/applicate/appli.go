package applicate

import (
	"IM_chat/dao/redisdao"
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"strconv"
)

func SearchAppli(user *models.AppliSearch) string {
	var err error
	if user.SendID, err = sql.SearchID(user.SendEmail); err != nil || user.UserID <= 0 {
		return errcode.Msg(errcode.NoPerson)
	}

	ok, err := redisdao.SetApplyLock(user.UserID, user.SendID)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if !ok {
		return errcode.Msg(errcode.ErrAlreadyRequested)
	}
	applyID, err := sql.SetApply(*user)
	if err != nil {
		return errcode.Msg(errcode.NotSendApply)
	}

	_ = redisdao.PushApplyInbox(user.SendID, applyID)
	_ = redisdao.PushApplySent(user.UserID, applyID)

	return errcode.Msg(errcode.SUCCESS)
}

func ListApp(userID int64) ([]models.Apply, string) {
	rows, err := sql.ListApplySent(userID)
	if err != nil {
		return nil, errcode.Msg(errcode.ERROR)
	}
	idSet := make(map[int64]struct{})
	list := make([]models.Apply, 0, len(rows)+16)

	for _, r := range rows {
		idSet[r.ID] = struct{}{}
		email, _ := sql.SearchEmail(r.ToID)
		name, _ := sql.SearchName(r.ToID)
		list = append(list, models.Apply{
			SendID:    r.ToID,
			SendEmail: email,
			SendName:  name,
			Msg:       r.Remark,
			Time:      r.CreateTime,
			Status:    r.Status,
		})
	}

	redisIDs, err := redisdao.GetApplySend(userID)
	if err == nil {
		for _, s := range redisIDs {
			id, _ := strconv.ParseInt(s, 10, 64)
			if id <= 0 {
				continue
			}
			if _, ok := idSet[id]; ok {
				continue
			}
			row, err := sql.GetApplyByID(id, userID)
			if err != nil {
				continue
			}
			idSet[id] = struct{}{}
			email, _ := sql.SearchEmail(row.ToID)
			name, _ := sql.SearchName(row.ToID)
			list = append(list, models.Apply{
				SendID:    row.ToID,
				SendEmail: email,
				SendName:  name,
				Msg:       row.Remark,
				Time:      row.CreateTime,
				Status:    row.Status,
			})
		}
	}

	return list, errcode.Msg(errcode.SUCCESS)
}

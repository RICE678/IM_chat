package applicate

import (
	"IM_chat/dao/redisdao"
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"strconv"
)

func SearchAppli(user *models.FindPerson) (find *models.FindMiddle, err1 string) {
	var err error
	if user.SendID, err = sql.SearchID(user.SendEmail); err != nil || user.UserID <= 0 {
		err1 = errcode.Msg(errcode.NoPerson)
		return
	}
	if user.SendID == user.UserID {
		err1 = errcode.Msg(errcode.NotAddMy)
		return
	}
	find = &models.FindMiddle{
		SendEmail: user.SendEmail,
		SendID:    user.SendID,
	}
	find.SendName, _ = sql.SearchName(user.SendID)
	find.SendPictures, _ = sql.SearchPicture(user.SendID)
	err1 = errcode.Msg(errcode.SUCCESS)
	return
}

func GetAppli(user *models.AppliSearch) string {
	ok, err := redisdao.SetApplyLock(user.UserID, user.SendID)
	if err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if !ok {
		return errcode.Msg(errcode.ErrAlreadyRequested)
	}
	user.ApplyID, err = sql.SetApply(*user)
	if err != nil {
		return errcode.Msg(errcode.NotSendApply)
	}
	user.Status = 0
	err = redisdao.SetApplyCache(user.ApplyID, user.UserID, user.SendID, user.Msg)
	if err != nil {
		return err.Error()
	}
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
		email_from, _ := sql.SearchEmail(r.FromID)
		email, _ := sql.SearchEmail(r.ToID)
		name, _ := sql.SearchName(r.ToID)
		picture, _ := sql.SearchPicture(r.ToID)
		list = append(list, models.Apply{
			FromID:      userID,
			SendID:      r.ToID,
			FromEmail:   email_from,
			SendEmail:   email,
			SendName:    name,
			Msg:         r.Remark,
			Time:        r.CreateTime,
			Status:      r.Status,
			SendPicture: picture,
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
			email_from, _ := sql.SearchEmail(row.FromID)
			picture, _ := sql.SearchPicture(row.ToID)
			list = append(list, models.Apply{
				SendID:      row.ToID,
				SendEmail:   email,
				SendName:    name,
				FromEmail:   email_from,
				Msg:         row.Remark,
				Time:        row.CreateTime,
				Status:      row.Status,
				SendPicture: picture,
			})
		}
	}

	return list, errcode.Msg(errcode.SUCCESS)
}
func RefuseFriend(friend *models.RefuseFriend) string {
	var err error
	if friend.AppliID, err = redisdao.GetApplyPair(friend.UserID, friend.Account_id); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = sql.ChangeStatusByPair(friend.AppliID, friend.Account_id, friend.UserID, 2); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.SetApplyStatus(friend.AppliID, 2); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.DelApplyFromInbox(friend.UserID, friend.AppliID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.DelApplyPair(friend.UserID, friend.Account_id); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}
func AcceptFriend(user *models.AcceptFriend) string {
	var err error
	if user.AppliID, err = redisdao.GetApplyPair(user.UserID, user.Account_id); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if user.AppliID <= 0 {
		return errcode.Msg(errcode.NoSend)
	}
	if err = sql.ChangeStatusByPair(user.AppliID, user.Account_id, user.UserID, 1); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.SetApplyStatus(user.AppliID, 1); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.DelApplyFromInbox(user.UserID, user.AppliID); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = redisdao.DelApplyPair(user.UserID, user.Account_id); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	if err = sql.CreateContact(*user); err != nil {
		return errcode.Msg(errcode.ERROR)
	}
	return errcode.Msg(errcode.SUCCESS)
}

func ShowList(userID int64) ([]models.Apply, string) {
	rows, err := sql.ListApplyTo(userID)
	if err != nil {
		return nil, errcode.Msg(errcode.ERROR)
	}
	idSet := make(map[int64]struct{})
	list := make([]models.Apply, 0, len(rows)+16)

	for _, r := range rows {
		idSet[r.ID] = struct{}{}
		email_from, _ := sql.SearchEmail(r.FromID)
		name, _ := sql.SearchName(r.FromID)
		picture, _ := sql.SearchPicture(r.FromID)
		list = append(list, models.Apply{
			FromID:      r.FromID,
			FromEmail:   email_from,
			SendID:      r.FromID,
			SendEmail:   email_from,
			SendName:    name,
			Msg:         r.Remark,
			Time:        r.CreateTime,
			Status:      r.Status,
			SendPicture: picture,
		})
	}

	redisIDs, err := redisdao.GetApplyInboxIDs(userID, 200)
	if err == nil {
		for _, id := range redisIDs {
			if _, ok := idSet[id]; ok {
				continue
			}
			row, err := sql.GetApplyByIDTo(id, userID)
			if err != nil {
				continue
			}
			idSet[id] = struct{}{}
			name, _ := sql.SearchName(row.FromID)
			email_from, _ := sql.SearchEmail(row.FromID)
			picture, _ := sql.SearchPicture(row.FromID)
			list = append(list, models.Apply{
				FromID:      row.FromID,
				SendID:      row.FromID,
				SendEmail:   email_from,
				SendName:    name,
				Msg:         row.Remark,
				FromEmail:   email_from,
				Time:        row.CreateTime,
				Status:      row.Status,
				SendPicture: picture,
			})
		}
	}

	return list, errcode.Msg(errcode.SUCCESS)
}

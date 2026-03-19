package applicate

import (
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
)

func SearchAppli(user *models.AppliSearch) string {
	var err error
	if user.SendID, err = sql.SearchID(user.SendEmail); err != nil || user.UserID <= 0 {
		return errcode.Msg(errcode.NoPerson)
	}

}

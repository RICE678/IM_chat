package userserve

import (
	"IM_chat/dao/redisdao"
	"IM_chat/pkg/errcode"
)

// HeartbeatHandler 服务端对应handler
func HeartbeatHandler(userID int64) string {
	if err := redisdao.RefreshOnline(userID); err != nil {
		return err.Error()
	}
	return errcode.Msg(errcode.SUCCESS)
}

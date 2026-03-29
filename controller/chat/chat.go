package chat

import (
	"IM_chat/middlewares"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (ChatController) SendMsg(c *gin.Context) {
	userVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
	}
	userID, ok := userVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}

}

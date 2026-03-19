package application

import (
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"IM_chat/service/applicate"
	"github.com/gin-gonic/gin"
)

type AppliController struct{}

func (AppliController) CreateAppli(c *gin.Context) {
	var req models.AppliSearch
	userIDVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	req.UserID, ok = userIDVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
	}
	res := applicate.SearchAppli(&req)
	c.JSON(200, res)
	return
}

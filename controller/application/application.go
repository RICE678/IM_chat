package application

import (
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"IM_chat/service/applicate"
	"github.com/gin-gonic/gin"
)

type AppliController struct{}

func NewAppliController() *AppliController {
	return new(AppliController)
}

// CreateAppli godoc
// @Summary 申请请求添加好友
// @Tags application
// @Accept json
// @Produce json
// @Param request body models.AppliSearch true "申请请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /a
// @Router /application/create [post]
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
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := applicate.SearchAppli(&req)
	c.JSON(200, res)
}

// ListAppli godoc
// @Summary 查看全部已发申请
// @Description 返回我发出的好友申请列表（合并 Redis + MySQL，含过期、成功、失败等全部状态）
// @Tags application
// @Accept json
// @Produce json
// @Success 200 {object} models.ListAppResponse "msg=success 时 list 为申请数组"
// @Security BearerAuth
// @Router /application/list [get]
func (AppliController) ListAppli(c *gin.Context) {
	userIDVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}
	list, msg := applicate.ListApp(userID)
	if list == nil {
		c.JSON(200, gin.H{"msg": msg, "list": []models.Apply{}})
		return
	}
	c.JSON(200, gin.H{"msg": msg, "list": list})
}

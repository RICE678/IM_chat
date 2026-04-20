package application

import (
	"IM_chat/logic/applicate"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type AppliController struct{}

func NewAppliController() *AppliController {
	return new(AppliController)
}

// SearchAppli godoc
// @Summary 搜索邮箱查看待添加好友
// @Tags application
// @Produce json
// @Param request body models.FindPerson true "按邮箱搜索"
// @Success 200 {object} models.FindEnd
// @Security BearerAuth
// @Router /application/search/email [post]
func (AppliController) SearchAppli(c *gin.Context) {
	userIDVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	var req models.FindPerson
	req.UserID, ok = userIDVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res, err := applicate.SearchAppli(&req)
	c.JSON(200, models.FindEnd{
		Err:  err,
		Find: res,
	})
}

// SearchNameAppli godoc
// @Summary 搜索名字查看待添加好友
// @Tags application
// @Accept json
// @Produce json
// @Param request body models.FindNamePerson true "按名字模糊搜索"
// @Success 200 {object} models.FindNameEnd
// @Security BearerAuth
// @Router /application/search/name [post]
func (AppliController) SearchNameAppli(c *gin.Context) {
	userIDVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	var req models.FindNamePerson
	req.UserID, ok = userIDVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res, err := applicate.SearchNameAppli(&req)
	c.JSON(200, models.FindNameEnd{
		Err:  err,
		Find: res,
	})
}

// CreateAppli godoc
// @Summary 申请请求添加好友
// @Tags application
// @Accept json
// @Produce json
// @Param request body models.AppliSearch true "申请请求"
// @Success 200 {string} string
// @Security BearerAuth
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
	res := applicate.GetAppli(&req)
	c.JSON(200, res)
}

// ListMyAppli godoc
// @Summary 查看全部已发申请
// @Description 返回我发出的好友申请列表
// @Tags application
// @Accept json
// @Produce json
// @Success 200 {object} models.ListAppResponse "msg=success 时 list 为申请数组"
// @Security BearerAuth
// @Router /application/mylist [get]
func (AppliController) ListMyAppli(c *gin.Context) {
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

// RefuseAppli godoc
// @Summary 拒绝好友申请
// @Tags application
// @Accept json
// @Produce json
// @Param request body models.RefuseFriend true "拒绝申请请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /application/refuse [put]
func (AppliController) RefuseAppli(c *gin.Context) {
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
	var req models.RefuseFriend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"msg": "ParamErr", "err": err.Error()})
		return
	}
	req.UserID = userID
	res := applicate.RefuseFriend(&req)
	c.JSON(200, res)
}

// AcceptAppli godoc
// @Summary 同意好友申请
// @Tags application
// @Accept json
// @Produce json
// @Param request body models.AcceptFriend true "同意申请请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /application/accept [put]
func (AppliController) AcceptAppli(c *gin.Context) {
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
	var req models.AcceptFriend
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"msg": "ParamErr", "err": err.Error()})
		return
	}
	req.UserID = userID
	res := applicate.AcceptFriend(&req)
	c.JSON(200, res)
}

// ListAppli godoc
// @Summary 查看全部接收到的申请
// @Description 返回我收到的好友申请列表
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
	list, msg := applicate.ShowList(userID)
	if list == nil {
		c.JSON(200, gin.H{"msg": msg, "list": []models.Apply{}})
		return
	}
	c.JSON(200, gin.H{"msg": msg, "list": list})
}

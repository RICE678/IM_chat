package user

import (
	"IM_chat/models"
	user2 "IM_chat/service/userserve"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

// Register godoc
// @Summary 用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.RegisterServer true "注册请求"
// @Success 200 {string} string
// @Router /user/register [post]
func Register(c *gin.Context) {
	var req models.RegisterServer
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := user2.RegisterDetail(&req)
	c.JSON(200, res)
}

// Login godoc
// @Summary 用户登录
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.LoginServer true "登录请求"
// @Success 200 {string} string
// @Router /user/login [post]
func Login(c *gin.Context) {
	var req models.LoginServer
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := user2.LoginDetail(&req)
	c.JSON(200, res)
}

// ReUpdate godoc
// @Summary 修改密码
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.ReUpdate true "修改密码请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/update/pwd [put]
func ReUpdate(c *gin.Context) {
	var req models.ReUpdate
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := user2.UpdateDetail(&req)
	c.JSON(200, res)
}

// ReEmail godoc
// @Summary 修改邮箱
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.ReEmail true "修改邮箱请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/update/email [put]
func ReEmail(c *gin.Context) {
	var req models.ReEmail
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
	}
	res := user2.ReEmail(&req)
	c.JSON(200, res)
}

// Heartbeat godoc
// @Summary 心跳续期
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.HeartbeatRequest true "心跳请求"
// @Success 200 {string} string
// @Router /user/heartbeat [post]
func Heartbeat(c *gin.Context) {
	var req models.HeartbeatRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := user2.HeartbeatHandler(req.UserID)
	c.JSON(200, res)
}

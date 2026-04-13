package user

import (
	user2 "IM_chat/logic/userserve"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

// Register godoc
// @Summary 用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.RegisterServer true "注册请求"
// @Success 200 {string} string
// @Router /user/register [post]
func (uc UserController) Register(c *gin.Context) {
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
// @Success 200 {object} models.LoginResponse
// @Router /user/login [post]
func (uc UserController) Login(c *gin.Context) {
	var req models.LoginServer
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	msg, token, userID := user2.LoginDetail(&req)
	c.JSON(200, models.LoginResponse{
		Msg:    msg,
		Token:  token,
		UserID: userID,
	})
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
func (uc UserController) ReUpdate(c *gin.Context) {
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
	var req models.ReUpdate
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	req.UserID = userID
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
func (uc UserController) ReEmail(c *gin.Context) {
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
	var req models.ReEmail
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	req.UserID = userID
	res := user2.ReEmail(&req)
	c.JSON(200, res)
}

// CreateUserMain godoc
// @Summary 完善用户资料
// @Tags user
// @Accept json
// @Produce json
// @Param request body models.UserMain2 true "完善资料：name 空则保留原名；gender/signature/picture_id 省略则不改；signature 可传空串清空；picture_id 为图片库 id（≥1），<1 时按默认头像 id=1 保存"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/create [post]
func (UserController) CreateUserMain(c *gin.Context) {
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
	var req models.UserMain2
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr: "+err.Error())
		return
	}
	req.UserID = userID
	res := user2.CreateAccountDetails(&req)
	c.JSON(200, res)
}

// ShowPictures godoc
// @Summary 展示图片库
// @Description 返回 picture 表全部 id、web（需登录，无请求体）
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.PictureMainReturn
// @Failure 401 {string} string "未登录或 token 无效"
// @Security BearerAuth
// @Router /user/show/pictures [post]
func (UserController) ShowPictures(c *gin.Context) {
	_, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	res, err := user2.SearchPictures()
	c.JSON(200, models.PictureMainReturn{
		User: res,
		Err:  err,
	})
	return
}

// LookUserMain godoc
// @Summary 显示当前用户资料
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.UserMainReturn
// @Security BearerAuth
// @Router /user/show/main [post]
func (UserController) LookUserMain(c *gin.Context) {
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
	res, err := user2.SearchAccountDetails(userID)
	c.JSON(200, &models.UserMainReturn{
		User: res,
		Err:  err,
	})

}

// Heartbeat godoc
// @Summary 心跳续期
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/heartbeat [post]
func (uc UserController) Heartbeat(c *gin.Context) {
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
	res := user2.HeartbeatHandler(userID)
	c.JSON(200, res)
}

// DelUser godoc
// @Summary 删除用户
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/deleteUser [delete]
func (uc UserController) DelUser(c *gin.Context) {
	userIDVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
	}
	res := user2.DelUserDetail(userID)
	c.JSON(200, res)
	return
}

// ReCode godoc
// @Summary 发送验证码以便验证身份
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Security BearerAuth
// @Router /user/pwd/code/send [post]
func (UserController) ReCode(c *gin.Context) {
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
	res := user2.ReCodeSend(userID)
	c.JSON(200, res)
	return
}

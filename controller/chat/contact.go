package chat

import (
	"IM_chat/logic/contact"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type ContactController struct{}

func NewContactController() *ContactController { return new(ContactController) }

// SearchContact godoc
// @Summary 查看通讯录
// @Tags Contact
// @Produce JSON
// @Success 200 {object} models.Contact
// @Security BearerAuth
// @Router /contact/list [post]
func (ContactController) SearchContact(c *gin.Context) {
	userVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	userID, ok := userVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}
	res, err := contact.SearchContactList(userID)
	c.JSON(200, models.Contact{
		List: res,
		Err:  err,
	})
	return
}

// SearchContactMain godoc
// @Summary 查看朋友详细信息
// @Tags Contact
// @Accept JSON
// @Produce JSON
// @Param request body models.FriendMain true "查看请求"
// @Success 200 {object} models.ContactMain
// @Security BearerAuth
// @Router /contact/friend/main [post]
func (ContactController) SearchContactMain(c *gin.Context) {
	userVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	var req models.FriendMain
	userID, ok := userVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}
	req.UserID = userID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res, err := contact.SearchFriendAll(req.FriendID, req.UserID)
	if err != errcode.Msg(errcode.SUCCESS) {
		c.JSON(200, gin.H{
			"list": nil,
			"err":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"list": res,
		"err":  err,
	})
}

// ChangeRemark godoc
// @Summary 修改朋友备注
// @Tags Contact
// @Produce JSON
// @Param request body models.FriendRemark true "修改备注请求"
// @Success 200 {STRING} string
// @Security BearerAuth
// @Router /contact/change/remark [post]
func (ContactController) ChangeRemark(c *gin.Context) {
	userVal, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeNeedLogin))
		return
	}
	var req models.FriendRemark
	userID, ok := userVal.(int64)
	if !ok {
		c.JSON(401, errcode.Msg(errcode.CodeInvalidToken))
		return
	}
	req.UserID = userID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := contact.RemarkChange(req.UserID, req.FriendID, req.Remark)
	c.JSON(200, res)
}

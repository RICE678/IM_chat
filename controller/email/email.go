package email

import (
	"IM_chat/service/email"
	"github.com/gin-gonic/gin"
)

type EmailController struct{}

func NewEmailController() EmailController {
	return EmailController{}
}

// ConfirmUserEmail godoc
// @Summary 发送邮箱验证码
// @Tags email
// @Accept json
// @Produce json
// @Param request body email.UserSendConfirmEmailService true "邮箱验证码请求"
// @Success 200 {string} string
// @Router /email/send [post]
func (ec EmailController) ConfirmUserEmail(c *gin.Context) {
	var req email.UserSendConfirmEmailService
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := req.SendConfirmEmail()
	c.JSON(200, res)
}

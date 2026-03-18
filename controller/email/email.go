package email

import (
	"IM_chat/service/email"
	"github.com/gin-gonic/gin"
)

// ConfirmUserEmail godoc
// @Summary 发送邮箱验证码
// @Tags email
// @Accept json
// @Produce json
// @Param request body email.UserSendConfirmEmailService true "邮箱验证码请求"
// @Success 200 {string} string
// @Router /email/send [post]
func ConfirmUserEmail(c *gin.Context) {
	var service email.UserSendConfirmEmailService
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	res := service.SendConfirmEmail()
	c.JSON(200, res)
}

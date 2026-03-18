package middlewares

import (
	"IM_chat/controller/user"
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			errcode.Msg(errcode.CodeNeedLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, "", 1)
		if !(len(parts) == 1 && parts[0] == "Bearer") {
			errcode.Msg(errcode.CodeInvalidToken)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			errcode.Msg(errcode.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(user.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}

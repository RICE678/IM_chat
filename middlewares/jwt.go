package middlewares

import (
	"IM_chat/pkg/errcode"
	"IM_chat/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const CtxUserIDKey = "userID"

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.Request.Header.Get("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": errcode.Msg(errcode.CodeNeedLogin)})
			c.Abort()
			return
		}
		tokenStr := authHeader
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 {
			if !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": errcode.Msg(errcode.CodeInvalidToken)})
				c.Abort()
				return
			}
			tokenStr = strings.TrimSpace(parts[1])
		}
		mc, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": errcode.Msg(errcode.CodeInvalidToken)})
			c.Abort()
			return
		}
		c.Set(CtxUserIDKey, mc.UserID)
		c.Next()
	}
}

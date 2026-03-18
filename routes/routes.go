package routes

import (
	"IM_chat/controller/email"
	"IM_chat/controller/user"
	"IM_chat/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// Setup 初始化路由，挂载所有接口
func Setup() *gin.Engine {
	r := gin.New()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	emailGroup := r.Group("/email")
	{
		emailGroup.POST("/send", email.ConfirmUserEmail)
	}
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", user.Register)
		userGroup.POST("/login", user.Login)
		userGroup.POST("/heartbeat", user.Heartbeat)
		userGroup.Use(middlewares.JWTAuthMiddleware())
		{
			userGroup.PUT("/update/pwd", user.ReUpdate)
			userGroup.PUT("/update/email", user.ReEmail)
		}
	}

	return r
}

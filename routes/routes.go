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
	emailController := email.NewEmailController()
	userController := user.NewUserController()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	emailGroup := r.Group("/email")
	{
		emailGroup.POST("/send", emailController.ConfirmUserEmail)
	}
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
		userGroup.Use(middlewares.JWTAuthMiddleware())
		{
			userGroup.POST("/heartbeat", userController.Heartbeat)
			userGroup.PUT("/update/pwd", userController.ReUpdate)
			userGroup.PUT("/update/email", userController.ReEmail)
			userGroup.POST("/create", userController.CreateUserMain)
			userGroup.DELETE("/deleteUser", userController.DelUser)
		}
	}
	applicationGroup := r.Group("/application")
	{
		applicationGroup.Use(middlewares.JWTAuthMiddleware())
		{

		}
	}
	return r
}

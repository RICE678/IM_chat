package routes

import (
	"IM_chat/controller/application"
	"IM_chat/controller/chat"
	"IM_chat/controller/email"
	"IM_chat/controller/user"
	"IM_chat/logic/chatupload"
	"IM_chat/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 32 << 20
	r.Static(chatupload.UploadURLPrefix, chatupload.DiskDir())
	emailController := email.NewEmailController()
	userController := user.NewUserController()
	applyController := application.NewAppliController()
	chatController := chat.NewChatController()
	contactController := chat.NewContactController()
	r.Use(middlewares.Cors(), middlewares.GinLogger(), middlewares.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFile("/docs/asyncapi.yaml", "./docs/asyncapi.yaml")
	r.StaticFile("/docs/socketio.md", "./docs/socketio_protocol.md")
	r.StaticFile("/docs/socketio", "./docs/asyncapi.html")
	r.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs/socketio")
	})
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
			userGroup.POST("/pwd/code/send", userController.ReCode)
			userGroup.POST("/heartbeat", userController.Heartbeat)
			userGroup.PUT("/update/pwd", userController.ReUpdate)
			userGroup.PUT("/update/email", userController.ReEmail)
			userGroup.POST("/show/main", userController.LookUserMain)
			userGroup.POST("/create", userController.CreateUserMain)
			userGroup.POST("/show/pictures", userController.ShowPictures)
			userGroup.DELETE("/deleteUser", userController.DelUser)
		}
	}
	applicationGroup := r.Group("/application")
	{
		applicationGroup.Use(middlewares.JWTAuthMiddleware())
		{
			applicationGroup.POST("/search/email", applyController.SearchAppli)
			applicationGroup.POST("/search/name", applyController.SearchNameAppli)
			applicationGroup.POST("/create", applyController.CreateAppli)
			applicationGroup.GET("/mylist", applyController.ListMyAppli)
			applicationGroup.PUT("/refuse", applyController.RefuseAppli)
			applicationGroup.PUT("/accept", applyController.AcceptAppli)
			applicationGroup.GET("/list", applyController.ListAppli)
		}
	}
	chatGroup := r.Group("/chat")
	{
		chatGroup.Use(middlewares.JWTAuthMiddleware())
		{
			chatGroup.GET("/show/all", chatController.ShowFriend)
			chatGroup.POST("/del", chatController.FriendChatDel)
			chatGroup.GET("/history", chatController.SearchHistory)
			chatGroup.POST("/read", chatController.EnterRead)
			chatGroup.POST("/upload", chatController.UploadChatFile)
		}
	}
	contactGroup := r.Group("/contact")
	{
		contactGroup.Use(middlewares.JWTAuthMiddleware())
		{
			contactGroup.POST("/list", contactController.SearchContact)
			contactGroup.POST("/friend/main", contactController.SearchContactMain)
			contactGroup.POST("/change/remark", contactController.ChangeRemark)
		}
	}
	r.Any("/socket.io/*any", chatController.ServeSocketIO)
	return r
}

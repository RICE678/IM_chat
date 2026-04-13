package chat

import (
	"IM_chat/logic/ws"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ChatController struct {
	SocketServer *socketio.Server
}

func NewChatController() *ChatController {
	server := ws.NewSocketIOServer()
	go server.Serve()
	return &ChatController{SocketServer: server}
}

func (cc *ChatController) ServeSocketIO(c *gin.Context) {
	cc.SocketServer.ServeHTTP(c.Writer, c.Request)
}

func (cc *ChatController) Shutdown() {
	_ = cc.SocketServer.Close()
}

func (cc *ChatController) SendMsg(c *gin.Context) {
	cc.ServeSocketIO(c)
}

// SearchHistory godoc
// @Summary 查看历史聊天记录
// @Tags chat
// @Accept json
// @Produce json
// @Param receiver_id query int true "对方用户 ID"
// @Param page query int false "页码，默认 1"
// @Param size query int false "每页条数，默认 20，最大 100"
// @Success 200 {object} models.HistoryResponse
// @Security BearerAuth
// @Router /chat/history [get]
func (ChatController) SearchHistory(c *gin.Context) {
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
	var req models.HistoryMsg
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 20
	}
	req.SenderID = userID
	err := ws.HistoryMain(&req)
	c.JSON(200, models.HistoryResponse{
		Err: err,
		Msg: &req,
	})
}

// SearchUnread godoc
// @Summary 私聊
// @Description 展示侧边栏私聊情况
// @Tags chat
// @Accept json
// @Produce json
// @Success 200 {object} models.ContactResponse
// @Security BearerAuth
// @Router /chat/unread [get]
func (ChatController) SearchUnread(c *gin.Context) {
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
	friends, err := ws.UnReadMain(userID)
	c.JSON(200, models.ContactResponse{
		Err:    err,
		Friend: friends,
	})

}

// EnterRead godoc
// @Summary 更改消息为已读
// @Tags chat
// @Accept json
// @Produce json
// @Param request body models.ReadContact true "会话已读请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /chat/read [post]
func (ChatController) EnterRead(c *gin.Context) {
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
	var req models.ReadContact
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	req.UserID = userID
	err := ws.ReadMain(&req)
	c.JSON(200, err)
}

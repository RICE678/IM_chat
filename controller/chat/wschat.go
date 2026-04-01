package chat

import (
	"IM_chat/logic/ws"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg godoc
// @Summary 私聊 WebSocket 收发消息
// @Description 建立私聊长连接。鉴权支持 Authorization: Bearer <token>，也支持 query token
// @Tags chat
// @Security BearerAuth
// @Param token query string false "JWT token"
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} map[string]string "未登录或 token 无效"
// @Router /chat/pm [get]
func (ChatController) SendMsg(c *gin.Context) {
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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, "websocket upgrade failed")
		return
	}
	client := ws.NewClient(userID, conn)
	ws.GlobalManager.Register(client)
	go client.WritePump()
	go client.ReadPump()
}

// SearchHistory godoc
// @Summary 查看历史聊天记录
// @Tags chat
// @Accept json
// @Produce json
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

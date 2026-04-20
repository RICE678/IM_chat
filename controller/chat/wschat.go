package chat

import (
	"IM_chat/logic/ws"
	"IM_chat/middlewares"
	"IM_chat/models"
	"IM_chat/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	Handler http.Handler
}

func NewChatController() *ChatController {
	_, handler := ws.NewSocketIOServer()
	return &ChatController{Handler: handler}
}

// ServeSocketIO godoc
// @Summary Socket.IO 私聊通道
// @Description 通过 query 参数 token 建立 Socket.IO 连接，连接后客户端发送 `msg` 事件。
// @Description `msg_type=1` 表示文本消息：需要 `msg`。
// @Description `msg_type=2` 表示文件消息：建议同时传 `file_url`、`file_name`，并将 `msg` 设为 `file_url`。
// @Description 服务端会回推 `message` 事件，发送成功后会返回 `type=ack` 的确认消息。
// @Tags chat
// @Accept json
// @Produce json
// @Param token query string true "JWT Token"
// @Param any path string true "socket.io transport path"
// @Success 101 {string} string "Switching Protocols"
// @Router /socket.io/{any} [get]
func (cc *ChatController) ServeSocketIO(c *gin.Context) {
	cc.Handler.ServeHTTP(c.Writer, c.Request)
}

// ShowFriend godoc
// @Summary 展示所有未被删除的聊天框（侧边栏）
// @Tags chat
// @Produce json
// @Success 200 {object} models.Contact
// @Security BearerAuth
// @Router /chat/show/all [get]
func (ChatController) ShowFriend(c *gin.Context) {
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
	res, err := ws.SearchShowList(userID)
	c.JSON(200, models.Contact{
		List: res,
		Err:  err,
	})
	return
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

// EnterRead godoc
// @Summary 更改会话消息为已读
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

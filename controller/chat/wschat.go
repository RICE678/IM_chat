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
// @Summary      Socket.IO 私聊通道（非 HTTP 接口，详见下方协议说明）
// @Description  本接口不是普通的 HTTP 接口，而是 **Socket.IO 长连接入口**。
// @Description  Swagger/OpenAPI 2.0 无法完整描述事件流，前端请配合 `docs/socketio_protocol.md`
// @Description  或 `docs/asyncapi.yaml` 查看完整协议。下面给出最小必要说明：
// @Description
// @Description  ### 1. 建立连接
// @Description  - 路径：`/socket.io/`（默认 path，Socket.IO 客户端会自动拼接）
// @Description  - 鉴权：通过 query 参数携带登录后的 `token`（JWT）
// @Description  - 示例：`io("http://host:port", { query: { token: "<JWT>" }, transports: ["websocket"] })`
// @Description  - Token 为空或非法会被立即 `disconnect`
// @Description
// @Description  ### 2. 客户端 → 服务端 事件：`msg`
// @Description  payload 为一个 JSON 对象，结构见 `models.WsMsg`。关键字段：
// @Description  - `type` string：`private` 或 `text`，均表示私聊消息
// @Description  - `receiver_id` int64：接收方用户 ID（必填）
// @Description  - `msg_type` int：`1=文本消息`，`2=文件消息`
// @Description  - `msg` string：文本内容；当 `msg_type=2` 时，可把文件 URL 放这里
// @Description  - `file_url` string：文件地址（`msg_type=2` 建议必填；若只传 `msg`，服务端会自动把 `msg` 赋值给 `file_url`）
// @Description  - `file_name` string：文件名（`msg_type=2` 建议传）
// @Description  - `sender_id` 由服务端从 token 解析，**前端不需要传**，传了也会被覆盖
// @Description
// @Description  **文本消息示例：**
// @Description  ```json
// @Description  { "type": "private", "receiver_id": 1024, "msg_type": 1, "msg": "你好呀" }
// @Description  ```
// @Description  **文件消息示例：**
// @Description  ```json
// @Description  { "type": "private", "receiver_id": 1024, "msg_type": 2,
// @Description    "msg": "https://cdn.example.com/f/abc.png",
// @Description    "file_url": "https://cdn.example.com/f/abc.png",
// @Description    "file_name": "abc.png" }
// @Description  ```
// @Description
// @Description  ### 3. 服务端 → 客户端 事件：`message`
// @Description  收到消息或 ack 时，服务端会向对应用户 emit `message` 事件，payload 同样是 `models.WsMsg`。
// @Description  - 对**发送方**：收到 `type=ack` 的回执，代表服务端已受理（`msg="ok"`）
// @Description  - 对**接收方**：收到 `type=private` 的消息体，带 `sender_id` / `msg_type` / `msg` / `file_url` / `file_name`
// @Description
// @Description  ### 4. 断开连接
// @Description  直接 `socket.disconnect()` 即可，服务端会自动清理在线态。
// @Description  同一个用户重复登录时，旧连接会被服务端主动踢下线。
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        token query string true "登录后的 JWT Token"
// @Param        EIO   query string false "Socket.IO Engine.IO 版本，一般由客户端 SDK 自动填" default(4)
// @Param        transport query string false "传输方式，推荐 websocket" Enums(websocket, polling)
// @Param        any   path  string true "Socket.IO 握手/轮询路径，由客户端 SDK 自动补全"
// @Success      101 {string} string "WebSocket Switching Protocols（握手成功）"
// @Success      200 {object} models.WsMsg "服务端回推的 message 事件 payload（仅用于文档展示，实际通过 Socket.IO 事件传递）"
// @Failure      401 {string} string "token 缺失或非法，连接会被立即关闭"
// @Router       /socket.io/{any} [get]
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

// FriendChatDel godoc
// @Summary 删除好友聊天框
// @Tags Contact
// @Accept json
// @Produce json
// @Param request body models.DelResponse true "删除请求"
// @Success 200 {string} string
// @Security BearerAuth
// @Router /chat/del [post]
func (ChatController) FriendChatDel(c *gin.Context) {
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
	var req models.DelResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, "ParamErr")
		return
	}
	req.UserID = userID
	res := ws.DelFriendMain(req)
	c.JSON(200, res)
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

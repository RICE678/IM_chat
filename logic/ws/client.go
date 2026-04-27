package ws

import (
	"IM_chat/models"
	"encoding/json"
	"sync"

	"github.com/zishang520/socket.io/servers/socket/v3"
	"go.uber.org/zap"
)

type Client struct {
	UserID int64
	Sock   *socket.Socket
	once   sync.Once
}

func NewClient(userID int64, sock *socket.Socket) *Client {
	return &Client{
		UserID: userID,
		Sock:   sock,
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		GlobalManager.Unregister(c)
		c.Sock.Disconnect(true)
	})
}

func (c *Client) Send(msg *models.WsMsg) {
	data, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("marshal msg failed", zap.Error(err))
		return
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		zap.L().Error("unmarshal to map failed", zap.Error(err))
		return
	}
	if emitErr := c.Sock.Emit("message", raw); emitErr != nil {
		zap.L().Error("sock emit failed", zap.Int64("userID", c.UserID), zap.Error(emitErr))
		return
	}
	zap.L().Info("sock emit ok", zap.Int64("userID", c.UserID), zap.Any("payload", raw))
}

package ws

import (
	"IM_chat/models"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
	"sync"
)

type Client struct {
	UserID int64
	Conn   socketio.Conn
	once   sync.Once
}

func NewClient(userID int64, conn socketio.Conn) *Client {
	return &Client{
		UserID: userID,
		Conn:   conn,
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		GlobalManager.Unregister(c)
		_ = c.Conn.Close()
	})
}

func (c *Client) Send(msg *models.WsMsg) {
	data, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("marshal msg failed", zap.Error(err))
		return
	}
	c.Conn.Emit("message", data)
}

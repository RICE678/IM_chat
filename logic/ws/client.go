package ws

import (
	"IM_chat/models"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
)

type Client struct {
	UserID int64
	Send   chan *models.WsMsg
	Conn   *websocket.Conn
	mu     sync.Mutex
	once   sync.Once
}

func NewClient(userID int64, conn *websocket.Conn) *Client {
	return &Client{
		UserID: userID,
		Send:   make(chan *models.WsMsg, 256),
		Conn:   conn,
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		GlobalManager.Unregister(c)
		close(c.Send)
		_ = c.Conn.Close()
	})
}

func (c *Client) ReadPump() {
	defer c.Close()
	c.Conn.SetReadLimit(maxMsgSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Warn("ws read err", zap.Error(err))
			}
			break
		}
		var msg models.WsMsg
		if err = json.Unmarshal(raw, &msg); err != nil {
			zap.L().Error("ws unmarshal err", zap.Error(err))
			continue
		}
		msg.SenderID = c.UserID

		HandleMessage(c, &msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, _ := json.Marshal(msg)
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

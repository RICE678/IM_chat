package ws

import (
	"IM_chat/dao/sql"
	"IM_chat/models"
	"IM_chat/pkg/jwt"
	"encoding/json"
	"net/http"

	"github.com/zishang520/socket.io/servers/socket/v3"
	"go.uber.org/zap"
)

func NewSocketIOServer() (*socket.Server, http.Handler) {
	io := socket.NewServer(nil, nil)

	io.On("connection", func(args ...any) {
		sock := args[0].(*socket.Socket)

		token := sock.Handshake().Query.Query().Get("token")
		if token == "" {
			sock.Disconnect(true)
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			sock.Disconnect(true)
			return
		}

		client := NewClient(claims.UserID, sock)
		sock.SetData(client)
		GlobalManager.Register(client)
		PushUnreadMessages(client)

		sock.On("msg", func(datas ...any) {
			if len(datas) == 0 {
				return
			}
			raw, err := json.Marshal(datas[0])
			if err != nil {
				zap.L().Error("marshal event data failed", zap.Error(err))
				return
			}
			var msg models.WsMsg
			if err := json.Unmarshal(raw, &msg); err != nil {
				zap.L().Error("unmarshal WsMsg failed", zap.Error(err))
				return
			}
			msg.SenderID = client.UserID
			if msg.MsgType == 2 {
				if msg.Msg == "" && msg.FileURL != "" {
					msg.Msg = msg.FileURL
				}
				if msg.Msg == "" {
					zap.L().Warn("file msg missing url", zap.Int64("sender_id", msg.SenderID), zap.Int("msg_type", msg.MsgType))
					return
				}
				if msg.FileURL == "" {
					msg.FileURL = msg.Msg
				}
			}
			HandleMessage(client, &msg)
		})

		sock.On("disconnect", func(...any) {
			client.Close()
		})
	})

	return io, io.ServeHandler(nil)
}

func PushUnreadMessages(client *Client) {
	msgs, err := sql.GetUnreadMessages(client.UserID)

	if err != nil {
		zap.L().Error("get unread messages failed", zap.Int64("userID", client.UserID), zap.Error(err))
		return
	}
	for _, msg := range msgs {
		content := msg.Msg
		if content == "" {
			content = msg.Context
		}
		wsMsg := &models.WsMsg{
			Type:       "private",
			Msg:        content,
			SenderID:   msg.UserID,
			ReceiverID: msg.ReceiverID,
			MsgType:    msg.MsgType,
			Timestamp:  msg.Timestamp,
		}
		if wsMsg.MsgType == 2 {
			wsMsg.FileURL = wsMsg.Msg
		}
		client.Send(wsMsg)
	}
	zap.L().Info("push unread messages done", zap.Int64("userID", client.UserID), zap.Int("count", len(msgs)))
}

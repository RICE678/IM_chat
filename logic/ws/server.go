package ws

import (
	"IM_chat/models"
	"IM_chat/pkg/jwt"
	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
)

func NewSocketIOServer() *socketio.Server {
	server := socketio.NewServer(nil)
	server.OnConnect("/", func(conn socketio.Conn) error {
		u := conn.URL()
		token := u.Query().Get("token")
		if token == "" {
			_ = conn.Close()
			return nil
		}
		userID, err := jwt.ParseToken(token)
		if err != nil {
			_ = conn.Close()
			return nil
		}
		ID := userID.UserID
		client := NewClient(ID, conn)
		conn.SetContext(client)
		GlobalManager.Register(client)
		return nil
	})
	server.OnEvent("/", "msg", func(conn socketio.Conn, msg models.WsMsg) {
		client, ok := conn.Context().(*Client)
		if !ok {
			return
		}
		msg.SenderID = client.UserID
		HandleMessage(client, &msg)
	})
	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		if client, ok := conn.Context().(*Client); ok {
			client.Close()
		}
	})
	server.OnError("/", func(conn socketio.Conn, err error) {
		zap.L().Error("socketio error", zap.Error(err))
	})
	return server
}

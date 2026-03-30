package ws

import (
	"IM_chat/models"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	clients map[int64]*Client
	mu      sync.Mutex
}

//var GlobalManager *Manager {
//clients:make(map[int64]*Client),
//}

func (m *Manager) Register(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if old, ok := m.clients[c.UserID]; ok {
		close(old.Send)
	}
	m.clients[c.UserID] = c
	zap.L().Info("user onlines", zap.Int64("userID", c.UserID))
}

func (m *Manager) Unregister(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if old, ok := m.clients[c.UserID]; ok && old == c {
		delete(m.clients, c.UserID)
		zap.L().Info("user offline", zap.Int64("userID", c.UserID))
	}
}

func (m *Manager) Send(userID int64, msg *models.ChatMsg) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if client, ok := m.clients[userID]; ok {
		select {
		case client.Send <- msg:
			return true
		default:
			return false
		}
	}
	return false
}

func (m *Manager) IsOnline(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.clients[userID]
	return ok
}

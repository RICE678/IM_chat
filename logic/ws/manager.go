package ws

import (
	"IM_chat/models"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	clients map[int64]*Client
	mu      sync.RWMutex
}

var GlobalManager = &Manager{
	clients: make(map[int64]*Client),
}

func (m *Manager) Register(c *Client) {
	m.mu.Lock()
	var old *Client
	if exist, ok := m.clients[c.UserID]; ok {
		old = exist
	}
	m.clients[c.UserID] = c
	m.mu.Unlock()
	if old != nil {
		old.Close()
	}
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

func (m *Manager) Send(userID int64, msg *models.WsMsg) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
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
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.clients[userID]
	return ok
}

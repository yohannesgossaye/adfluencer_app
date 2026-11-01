package realtime

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"sync"

	"github.com/coder/websocket"
)

// Hub tracks open connections per user in this process. Running several
// backend instances needs a shared fan-out (e.g. Redis pub/sub) in front of it.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]bool)
	}
	h.clients[client.UserID][client] = true
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, ok := h.clients[client.UserID]
	if !ok {
		return
	}
	if _, exists := userClients[client]; exists {
		delete(userClients, client)
		close(client.Send)
	}
	if len(userClients) == 0 {
		delete(h.clients, client.UserID)
	}
}

// SendToUser delivers an event to every connection of the user. A client
// whose buffer is full is closed so it reconnects and reloads, rather than
// silently missing events.
func (h *Hub) SendToUser(userID string, msg dto.WSOutgoingMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients[userID] {
		if !client.TrySend(msg) {
			client.Close(websocket.StatusPolicyViolation, "connection too slow, please reconnect")
		}
	}
}

// CloseAll tells every client the server is going away (close code 1001) so
// they reconnect promptly, e.g. during a deploy.
func (h *Hub) CloseAll() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userClients := range h.clients {
		for client := range userClients {
			client.Close(websocket.StatusGoingAway, "server restarting")
		}
	}
}

func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

package websocket

import (
    "encoding/json"
    "sync"
    "time"

    "user-api/internal/logger"

    "github.com/gofiber/contrib/websocket"
    "go.uber.org/zap"
)

type Event struct {
    Type      string    `json:"type"`
    UserID    int32     `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
}


type Hub struct {
    clients map[*websocket.Conn]bool
    mu      sync.RWMutex   // protects the clients map
}

func NewHub() *Hub {
    return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) Register(conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.clients[conn] = true
    logger.Log.Info("websocket client connected",
        zap.Int("total_clients", len(h.clients)))
}

func (h *Hub) Unregister(conn *websocket.Conn) {
    h.mu.Lock()
    defer h.mu.Unlock()
    delete(h.clients, conn)
    logger.Log.Info("websocket client disconnected",
        zap.Int("total_clients", len(h.clients)))
}

func (h *Hub) Broadcast(eventType string, userID int32) {
    event := Event{
        Type:      eventType,
        UserID:    userID,
        Timestamp: time.Now(),
    }
    data, err := json.Marshal(event)
    if err != nil {
        logger.Log.Error("ws broadcast marshal failed", zap.Error(err))
        return
    }

    h.mu.RLock()
    defer h.mu.RUnlock()

    for conn := range h.clients {
        if err := conn.WriteMessage(1, data); err != nil {
            logger.Log.Error("ws write failed", zap.Error(err))
        }
    }
}
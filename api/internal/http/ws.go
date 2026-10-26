package http

import (
	"encoding/json"
	"log"
	nethttp "net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type DashboardDetectionEvent struct {
	Type     string `json:"type"`
	DeviceID string `json:"device_id"`
	TS       int64  `json:"ts"`
}

type DetectionNotifier interface {
	Broadcast(DashboardDetectionEvent)
}

type wsClient struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu      sync.Mutex
	clients map[*wsClient]struct{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*nethttp.Request) bool { return true },
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*wsClient]struct{})}
}

func (h *Hub) Broadcast(event DashboardDetectionEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("marshal detection event: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
		}
	}
}

func (h *Hub) ServeWS(w nethttp.ResponseWriter, r *nethttp.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade websocket: %v", err)
		return
	}

	client := &wsClient{conn: conn, send: make(chan []byte, 16)}
	h.register(client)
	go client.writePump()
	client.readPump(h)
}

func (h *Hub) register(client *wsClient) {
	h.mu.Lock()
	h.clients[client] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) unregister(client *wsClient) {
	h.mu.Lock()
	delete(h.clients, client)
	h.mu.Unlock()
	close(client.send)
}

func (c *wsClient) writePump() {
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *wsClient) readPump(h *Hub) {
	defer func() {
		h.unregister(c)
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

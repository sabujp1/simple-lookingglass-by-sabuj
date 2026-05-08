package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents a WebSocket client
type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	UserID string
}

// Hub maintains active clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.ID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends message to all clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// SendToUser sends message to a specific user
func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  hub,
	}

	hub.register <- client

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse incoming message
		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Handle different message types
		c.handleMessage(&msg)
	}
}

// writePump pumps messages to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-c.Conn.CloseNotify():
			return
		}
	}
}

// handleMessage handles incoming WebSocket messages
func (c *Client) handleMessage(msg *WebSocketMessage) {
	switch msg.Type {
	case "auth":
		// Authenticate the connection
		c.UserID = msg.UserID
		c.sendJSON(WebSocketMessage{
			Type:    "auth_success",
			Payload: map[string]string{"client_id": c.ID},
		})

	case "query":
		// Start a query and stream results
		c.handleQuery(msg)

	case "ping":
		// Heartbeat
		c.sendJSON(WebSocketMessage{
			Type: "pong",
		})
	}
}

// handleQuery handles a query request
func (c *Client) handleQuery(msg *WebSocketMessage) {
	// Extract query parameters
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		c.sendJSON(WebSocketMessage{
			Type:    "error",
			Payload: map[string]string{"message": "invalid payload"},
		})
		return
	}

	queryID := uuid.New().String()

	// Send acknowledgment
	c.sendJSON(WebSocketMessage{
		Type:    "query_started",
		ID:      queryID,
		Payload: payload,
	})

	// In a real implementation, this would:
	// 1. Execute the query via the query service
	// 2. Stream results back via WebSocket
	// 3. Send completion message

	// Simulate query execution
	c.sendJSON(WebSocketMessage{
		Type:    "query_output",
		ID:      queryID,
		Payload: map[string]string{"output": "Query result would appear here..."},
	})

	c.sendJSON(WebSocketMessage{
		Type:    "query_completed",
		ID:      queryID,
		Payload: map[string]int{"duration_ms": 100},
	})
}

// sendJSON sends a JSON message
func (c *Client) sendJSON(msg WebSocketMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	c.Send <- data
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	ID      string      `json:"id,omitempty"`
	UserID  string      `json:"user_id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}
// Package websocket provides WebSocket handlers for real-time updates.
// It manages client connections and broadcasts pool/opportunity updates.
package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	MessageTypePoolUpdate       MessageType = "pool_update"
	MessageTypePoolsSnapshot    MessageType = "pools_snapshot"
	MessageTypeOpportunityAlert MessageType = "opportunity_alert"
	MessageTypePing             MessageType = "ping"
	MessageTypePong             MessageType = "pong"
	MessageTypeError            MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType     `json:"type"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID         string
	Conn       *websocket.Conn
	Send       chan []byte
	Hub        *Hub
	Subscribed map[string]bool // Subscribed channels
	mu         sync.RWMutex
}

// Hub manages WebSocket client connections and message broadcasting
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Channel-specific clients
	poolClients        map[*Client]bool
	opportunityClients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Configuration
	config config.WebSocketConfig

	mu sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub(cfg config.WebSocketConfig) *Hub {
	return &Hub{
		clients:            make(map[*Client]bool),
		poolClients:        make(map[*Client]bool),
		opportunityClients: make(map[*Client]bool),
		broadcast:          make(chan []byte, 256),
		register:           make(chan *Client),
		unregister:         make(chan *Client),
		config:             cfg,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Debug().Str("client_id", client.ID).Msg("Client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.poolClients, client)
				delete(h.opportunityClients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Debug().Str("client_id", client.ID).Msg("Client disconnected")

		case message := <-h.broadcast:
			h.mu.RLock()
			// Collect clients with full buffers
			var deadClients []*Client
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client buffer full, mark for removal
					deadClients = append(deadClients, client)
				}
			}
			h.mu.RUnlock()

			// Remove dead clients with write lock
			if len(deadClients) > 0 {
				h.mu.Lock()
				for _, client := range deadClients {
					if _, ok := h.clients[client]; ok {
						delete(h.clients, client)
						delete(h.poolClients, client)
						delete(h.opportunityClients, client)
						close(client.Send)
					}
				}
				h.mu.Unlock()
			}
		}
	}
}

// BroadcastPoolUpdate sends a pool update to all pool subscribers
func (h *Hub) BroadcastPoolUpdate(pool *models.Pool) {
	data, err := json.Marshal(pool)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal pool for broadcast")
		return
	}

	msg := Message{
		Type:      MessageTypePoolUpdate,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	h.mu.RLock()
	var deadClients []*Client
	for client := range h.poolClients {
		select {
		case client.Send <- msgBytes:
		default:
			// Client buffer full, mark for removal
			deadClients = append(deadClients, client)
		}
	}
	h.mu.RUnlock()

	// Clean up dead clients
	if len(deadClients) > 0 {
		h.mu.Lock()
		for _, client := range deadClients {
			delete(h.poolClients, client)
		}
		h.mu.Unlock()
	}
}

// BroadcastOpportunityAlert sends an opportunity alert to subscribers
func (h *Hub) BroadcastOpportunityAlert(opp *models.Opportunity) {
	data, err := json.Marshal(opp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal opportunity for broadcast")
		return
	}

	msg := Message{
		Type:      MessageTypeOpportunityAlert,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	h.mu.RLock()
	var deadClients []*Client
	for client := range h.opportunityClients {
		select {
		case client.Send <- msgBytes:
		default:
			// Client buffer full, mark for removal
			deadClients = append(deadClients, client)
		}
	}
	h.mu.RUnlock()

	// Clean up dead clients
	if len(deadClients) > 0 {
		h.mu.Lock()
		for _, client := range deadClients {
			delete(h.opportunityClients, client)
		}
		h.mu.Unlock()
	}
}

// SubscribeToPool adds a client to pool updates
func (h *Hub) SubscribeToPool(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.poolClients[client] = true
}

// SubscribeToOpportunities adds a client to opportunity alerts
func (h *Hub) SubscribeToOpportunities(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.opportunityClients[client] = true
}

// UnsubscribeFromPool removes a client from pool updates
func (h *Hub) UnsubscribeFromPool(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.poolClients, client)
}

// UnsubscribeFromOpportunities removes a client from opportunity alerts
func (h *Hub) UnsubscribeFromOpportunities(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.opportunityClients, client)
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return map[string]int{
		"total_clients":       len(h.clients),
		"pool_subscribers":    len(h.poolClients),
		"opp_subscribers":     len(h.opportunityClients),
	}
}

// NewClient creates a new WebSocket client
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:         id,
		Conn:       conn,
		Hub:        hub,
		Send:       make(chan []byte, 256),
		Subscribed: make(map[string]bool),
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(c.Hub.config.PingInterval)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// Channel closed
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Debug().Err(err).Str("client_id", c.ID).Msg("Write error")
				return
			}

		case <-ticker.C:
			// Send ping
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(c.Hub.config.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(c.Hub.config.PongTimeout))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(c.Hub.config.PongTimeout))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Debug().Err(err).Str("client_id", c.ID).Msg("WebSocket error")
			}
			break
		}

		// Handle incoming messages (subscriptions, etc.)
		c.handleMessage(message)
	}
}

// handleMessage processes incoming messages from clients
func (c *Client) handleMessage(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal client message")
		return
	}

	switch msg.Type {
	case MessageTypePing:
		// Respond with pong
		response := Message{
			Type:      MessageTypePong,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		responseBytes, _ := json.Marshal(response)
		c.Send <- responseBytes

	default:
		log.Debug().Str("type", string(msg.Type)).Msg("Received unknown message type")
	}
}

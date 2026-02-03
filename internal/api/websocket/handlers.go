package websocket

import (
	"context"
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
)

// Handler manages WebSocket connections
type Handler struct {
	hub       *Hub
	redisRepo *redis.Repository
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub, redisRepo *redis.Repository) *Handler {
	return &Handler{
		hub:       hub,
		redisRepo: redisRepo,
	}
}

// UpgradeCheck is middleware to check if the request is a WebSocket upgrade
func UpgradeCheck(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandlePoolUpdates handles WebSocket connections for pool updates
// WS /ws/pools
func (h *Handler) HandlePoolUpdates(c *websocket.Conn) {
	clientID := uuid.New().String()

	client := NewClient(clientID, c, h.hub)

	// Register client
	h.hub.register <- client

	// Subscribe to pool updates
	h.hub.SubscribeToPool(client)

	log.Info().
		Str("client_id", clientID).
		Str("remote_addr", c.RemoteAddr().String()).
		Msg("WebSocket client connected to pool updates")

	// Start read/write pumps
	go client.WritePump()
	client.ReadPump() // Blocking

	log.Info().
		Str("client_id", clientID).
		Msg("WebSocket client disconnected from pool updates")
}

// HandleOpportunityAlerts handles WebSocket connections for opportunity alerts
// WS /ws/opportunities
func (h *Handler) HandleOpportunityAlerts(c *websocket.Conn) {
	clientID := uuid.New().String()

	client := NewClient(clientID, c, h.hub)

	// Register client
	h.hub.register <- client

	// Subscribe to opportunity alerts
	h.hub.SubscribeToOpportunities(client)

	log.Info().
		Str("client_id", clientID).
		Str("remote_addr", c.RemoteAddr().String()).
		Msg("WebSocket client connected to opportunity alerts")

	// Start read/write pumps
	go client.WritePump()
	client.ReadPump() // Blocking

	log.Info().
		Str("client_id", clientID).
		Msg("WebSocket client disconnected from opportunity alerts")
}

// StartRedisSubscriber starts listening to Redis pub/sub channels
// and broadcasts messages to WebSocket clients
func (h *Handler) StartRedisSubscriber(ctx context.Context) {
	// Subscribe to pool updates
	go h.subscribeToPoolUpdates(ctx)

	// Subscribe to opportunity alerts
	go h.subscribeToOpportunityAlerts(ctx)
}

// subscribeToPoolUpdates listens to Redis pool update channel
func (h *Handler) subscribeToPoolUpdates(ctx context.Context) {
	pubsub := h.redisRepo.SubscribePoolUpdates(ctx)
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Info().Msg("Started Redis subscriber for pool updates")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping pool updates subscriber")
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Parse pool from message
			var pool models.Pool
			if err := json.Unmarshal([]byte(msg.Payload), &pool); err != nil {
				log.Debug().Err(err).Msg("Failed to unmarshal pool update")
				continue
			}

			// Broadcast to WebSocket clients
			h.hub.BroadcastPoolUpdate(&pool)
		}
	}
}

// subscribeToOpportunityAlerts listens to Redis opportunity channel
func (h *Handler) subscribeToOpportunityAlerts(ctx context.Context) {
	pubsub := h.redisRepo.SubscribeOpportunityAlerts(ctx)
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Info().Msg("Started Redis subscriber for opportunity alerts")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping opportunity alerts subscriber")
			return
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Parse opportunity from message
			var opp models.Opportunity
			if err := json.Unmarshal([]byte(msg.Payload), &opp); err != nil {
				log.Debug().Err(err).Msg("Failed to unmarshal opportunity alert")
				continue
			}

			// Broadcast to WebSocket clients
			h.hub.BroadcastOpportunityAlert(&opp)
		}
	}
}

// GetHubStats returns current WebSocket hub statistics
func (h *Handler) GetHubStats() map[string]int {
	return h.hub.GetStats()
}

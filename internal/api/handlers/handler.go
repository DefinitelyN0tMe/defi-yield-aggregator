// Package handlers contains HTTP request handlers for the DeFi Yield Aggregator API.
package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/elasticsearch"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/postgres"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
)

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	config *config.Config
	pg     *postgres.Repository
	redis  *redis.Repository
	es     *elasticsearch.Repository
	startTime time.Time
}

// NewHandler creates a new Handler with all dependencies
func NewHandler(
	cfg *config.Config,
	pg *postgres.Repository,
	redis *redis.Repository,
	es *elasticsearch.Repository,
) *Handler {
	return &Handler{
		config: cfg,
		pg:     pg,
		redis:  redis,
		es:     es,
		startTime: time.Now(),
	}
}

// HealthCheck returns the health status of the service and its dependencies
// GET /api/v1/health
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.Context()

	health := models.HealthCheck{
		Status:    "healthy",
		Version:   "1.0.0",
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  make(map[string]models.ServiceHealth),
	}

	// Check PostgreSQL
	pgStart := time.Now()
	pgErr := h.pg.Ping(ctx)
	health.Services["postgresql"] = models.ServiceHealth{
		Status:  boolToStatus(pgErr == nil),
		Latency: time.Since(pgStart).String(),
		Message: errToMessage(pgErr),
	}

	// Check Redis
	redisStart := time.Now()
	redisErr := h.redis.Ping(ctx)
	health.Services["redis"] = models.ServiceHealth{
		Status:  boolToStatus(redisErr == nil),
		Latency: time.Since(redisStart).String(),
		Message: errToMessage(redisErr),
	}

	// Check ElasticSearch
	esStart := time.Now()
	esErr := h.es.Ping(ctx)
	health.Services["elasticsearch"] = models.ServiceHealth{
		Status:  boolToStatus(esErr == nil),
		Latency: time.Since(esStart).String(),
		Message: errToMessage(esErr),
	}

	// Determine overall health
	if pgErr != nil || redisErr != nil {
		health.Status = "unhealthy"
		return c.Status(fiber.StatusServiceUnavailable).JSON(health)
	}
	if esErr != nil {
		health.Status = "degraded"
	}

	return c.JSON(health)
}

// ErrorHandler is the custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error code
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}

// Helper functions

func boolToStatus(ok bool) string {
	if ok {
		return "up"
	}
	return "down"
}

func errToMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

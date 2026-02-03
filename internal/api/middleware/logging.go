package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RequestLogger is a structured logging middleware for HTTP requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response status
		status := c.Response().StatusCode()

		// Build log event
		event := log.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}

		// Log request details
		event.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Int("status", status).
			Dur("latency", latency).
			Str("user_agent", c.Get("User-Agent")).
			Str("request_id", c.GetRespHeader("X-Request-ID")).
			Msg("HTTP Request")

		return err
	}
}

// RequestMetrics tracks request metrics for monitoring
type RequestMetrics struct {
	TotalRequests     int64
	SuccessRequests   int64
	ErrorRequests     int64
	TotalLatencyMs    int64
	RequestsByPath    map[string]int64
	RequestsByStatus  map[int]int64
}

var metrics = &RequestMetrics{
	RequestsByPath:   make(map[string]int64),
	RequestsByStatus: make(map[int]int64),
}

// MetricsCollector collects basic request metrics
func MetricsCollector() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Update metrics
		latency := time.Since(start)
		status := c.Response().StatusCode()
		path := c.Route().Path

		metrics.TotalRequests++
		metrics.TotalLatencyMs += latency.Milliseconds()
		metrics.RequestsByPath[path]++
		metrics.RequestsByStatus[status]++

		if status >= 400 {
			metrics.ErrorRequests++
		} else {
			metrics.SuccessRequests++
		}

		return err
	}
}

// GetMetrics returns current request metrics
func GetMetrics() *RequestMetrics {
	return metrics
}

// ResetMetrics resets all metrics counters
func ResetMetrics() {
	metrics = &RequestMetrics{
		RequestsByPath:   make(map[string]int64),
		RequestsByStatus: make(map[int]int64),
	}
}

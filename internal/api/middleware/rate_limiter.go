// Package middleware contains HTTP middleware for the DeFi Yield Aggregator API.
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
)

// RateLimiter creates a rate limiting middleware using a sliding window algorithm.
// It limits requests per IP address based on the configured thresholds.
func RateLimiter(cfg config.RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		// Maximum number of requests in the time window
		Max: cfg.Requests,

		// Time window for rate limiting
		Expiration: cfg.Window,

		// Use IP address as the key for rate limiting
		KeyGenerator: func(c *fiber.Ctx) string {
			// Try to get real IP from X-Forwarded-For header (for proxied requests)
			ip := c.Get("X-Forwarded-For")
			if ip == "" {
				ip = c.Get("X-Real-IP")
			}
			if ip == "" {
				ip = c.IP()
			}
			return ip
		},

		// Custom response when rate limit is exceeded
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    429,
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
		},

		// Skip rate limiting for health check endpoints
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			return path == "/health" || path == "/api/v1/health"
		},
	})
}

// SlowDown creates a middleware that adds artificial delay after threshold
// This is less aggressive than hard rate limiting
func SlowDown(threshold int, delay time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        threshold * 2, // Allow more requests but slow down
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			// Add delay instead of rejecting
			time.Sleep(delay)
			return c.Next()
		},
	})
}

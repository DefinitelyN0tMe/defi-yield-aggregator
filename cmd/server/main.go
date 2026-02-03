// Package main is the entry point for the DeFi Yield Aggregator API server.
// It initializes all dependencies and starts the HTTP/WebSocket server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/api/graphql"
	"github.com/maxjove/defi-yield-aggregator/internal/api/handlers"
	"github.com/maxjove/defi-yield-aggregator/internal/api/middleware"
	ws "github.com/maxjove/defi-yield-aggregator/internal/api/websocket"
	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/elasticsearch"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/postgres"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
)

// Build information - set via ldflags during build
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup structured logging
	setupLogger(cfg)

	log.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("git_commit", GitCommit).
		Str("environment", cfg.App.Env).
		Msg("Starting DeFi Yield Aggregator API Server")

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize PostgreSQL connection
	pgRepo, err := postgres.NewRepository(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pgRepo.Close()
	log.Info().Msg("Connected to PostgreSQL")

	// Initialize Redis connection
	redisRepo, err := redis.NewRepository(ctx, cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisRepo.Close()
	log.Info().Msg("Connected to Redis")

	// Initialize ElasticSearch connection
	esRepo, err := elasticsearch.NewRepository(cfg.ElasticSearch)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to ElasticSearch")
	}
	log.Info().Msg("Connected to ElasticSearch")

	// Create HTTP handler with dependencies
	h := handlers.NewHandler(cfg, pgRepo, redisRepo, esRepo)

	// Create WebSocket hub and handler
	wsHub := ws.NewHub(cfg.WebSocket)
	wsHandler := ws.NewHandler(wsHub, redisRepo)

	// Start WebSocket hub
	go wsHub.Run()
	log.Info().Msg("WebSocket hub started")

	// Start Redis subscriber for real-time updates
	go wsHandler.StartRedisSubscriber(ctx)

	// Create Fiber app with configuration
	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		DisableStartupMessage: cfg.IsProduction(),
		ErrorHandler:          handlers.ErrorHandler,
	})

	// Setup middleware
	setupMiddleware(app, cfg)

	// Create GraphQL resolver
	gqlResolver := graphql.NewResolver(pgRepo, redisRepo, esRepo)

	// Setup routes
	setupRoutes(app, h, wsHandler, gqlResolver)

	// Start server in goroutine
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	go func() {
		if err := app.Listen(serverAddr); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	log.Info().
		Str("address", serverAddr).
		Msg("Server started successfully")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Cancel context to stop background goroutines
	cancel()

	// Give outstanding requests 10 seconds to complete
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Error().Err(err).Msg("Error during server shutdown")
	}

	log.Info().Msg("Server stopped")
}

// setupLogger configures the zerolog logger based on environment
func setupLogger(cfg *config.Config) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	if cfg.IsDevelopment() {
		// Human-readable console output for development
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON output for production (better for log aggregation)
		zerolog.TimeFieldFormat = time.RFC3339Nano
	}
}

// setupMiddleware configures all middleware for the Fiber app
func setupMiddleware(app *fiber.App, cfg *config.Config) {
	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))

	// Request ID for tracing
	app.Use(requestid.New())

	// Request logging
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     stringSliceToString(cfg.CORS.AllowedOrigins),
		AllowMethods:     stringSliceToString(cfg.CORS.AllowedMethods),
		AllowHeaders:     stringSliceToString(cfg.CORS.AllowedHeaders),
		AllowCredentials: true,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	// Rate limiting (skip for WebSocket upgrades)
	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return middleware.RateLimiter(cfg.RateLimit)(c)
	})
}

// setupRoutes configures all API routes
func setupRoutes(app *fiber.App, h *handlers.Handler, wsHandler *ws.Handler, gqlResolver *graphql.Resolver) {
	// Health check (no versioning)
	app.Get("/health", h.HealthCheck)

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Health check (versioned)
	v1.Get("/health", h.HealthCheck)

	// Pool routes
	pools := v1.Group("/pools")
	pools.Get("/", h.ListPools)
	pools.Get("/:id", h.GetPool)
	pools.Get("/:id/history", h.GetPoolHistory)

	// Opportunity routes
	opportunities := v1.Group("/opportunities")
	opportunities.Get("/", h.ListOpportunities)
	opportunities.Get("/trending", h.GetTrendingPools)

	// Aggregated data routes
	v1.Get("/chains", h.ListChains)
	v1.Get("/protocols", h.ListProtocols)
	v1.Get("/stats", h.GetStats)

	// GraphQL routes
	app.Post("/graphql", gqlResolver.Handle)
	app.Get("/graphql", graphql.Playground) // GraphQL Playground UI

	// WebSocket routes
	wsGroup := app.Group("/ws")

	// WebSocket upgrade check middleware
	wsGroup.Use(ws.UpgradeCheck)

	// Pool updates WebSocket
	wsGroup.Get("/pools", websocket.New(wsHandler.HandlePoolUpdates, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))

	// Opportunity alerts WebSocket
	wsGroup.Get("/opportunities", websocket.New(wsHandler.HandleOpportunityAlerts, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}))

	// WebSocket stats endpoint (for monitoring)
	v1.Get("/ws/stats", func(c *fiber.Ctx) error {
		return c.JSON(wsHandler.GetHubStats())
	})
}

// Helper function to convert string slice to comma-separated string
func stringSliceToString(slice []string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}

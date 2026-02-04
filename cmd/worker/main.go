// Package main is the entry point for the DeFi Yield Aggregator background worker.
// It handles scheduled data fetching from external APIs and opportunity detection.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/elasticsearch"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/postgres"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
	"github.com/maxjove/defi-yield-aggregator/internal/services/analytics"
	"github.com/maxjove/defi-yield-aggregator/internal/services/coingecko"
	"github.com/maxjove/defi-yield-aggregator/internal/services/defillama"
	"github.com/maxjove/defi-yield-aggregator/internal/services/opportunity"
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
		Str("environment", cfg.App.Env).
		Msg("Starting DeFi Yield Aggregator Worker")

	// Initialize dependencies
	ctx := context.Background()

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

	// Create ElasticSearch indices
	if err := esRepo.CreateIndices(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to create ElasticSearch indices")
	}

	// Initialize API clients
	defiLlamaClient := defillama.NewClient(cfg.DeFiLlama)
	coinGeckoClient := coingecko.NewClient(cfg.CoinGecko)

	// Initialize services
	analyticsService := analytics.NewService(cfg.Scoring)
	opportunityService := opportunity.NewService(cfg.Worker, pgRepo, redisRepo, analyticsService)

	// Create scheduler
	scheduler := cron.New(cron.WithSeconds())

	// Schedule DeFiLlama fetch job (every 3 minutes)
	_, err = scheduler.AddFunc("0 */3 * * * *", func() {
		runDeFiLlamaJob(ctx, cfg, defiLlamaClient, pgRepo, redisRepo, esRepo, analyticsService)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to schedule DeFiLlama job")
	}
	log.Info().Str("interval", "3m").Msg("Scheduled DeFiLlama fetch job")

	// Schedule CoinGecko fetch job (every 10 minutes)
	_, err = scheduler.AddFunc("0 */10 * * * *", func() {
		runCoinGeckoJob(ctx, coinGeckoClient, redisRepo)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to schedule CoinGecko job")
	}
	log.Info().Str("interval", "10m").Msg("Scheduled CoinGecko fetch job")

	// Schedule opportunity detection job (every 5 minutes)
	_, err = scheduler.AddFunc("0 */5 * * * *", func() {
		runOpportunityDetectionJob(ctx, opportunityService, pgRepo, redisRepo)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to schedule opportunity detection job")
	}
	log.Info().Str("interval", "5m").Msg("Scheduled opportunity detection job")

	// Start scheduler
	scheduler.Start()
	log.Info().Msg("Worker scheduler started")

	// Run initial fetch immediately
	go func() {
		log.Info().Msg("Running initial data fetch...")
		runDeFiLlamaJob(ctx, cfg, defiLlamaClient, pgRepo, redisRepo, esRepo, analyticsService)
		runCoinGeckoJob(ctx, coinGeckoClient, redisRepo)
		runOpportunityDetectionJob(ctx, opportunityService, pgRepo, redisRepo)
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down worker...")

	// Stop scheduler gracefully
	stopCtx := scheduler.Stop()
	<-stopCtx.Done()

	log.Info().Msg("Worker stopped")
}

// setupLogger configures the zerolog logger based on environment
func setupLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		zerolog.TimeFieldFormat = time.RFC3339Nano
	}
}

// runDeFiLlamaJob fetches pools from DeFiLlama and stores them
func runDeFiLlamaJob(
	ctx context.Context,
	cfg *config.Config,
	client *defillama.Client,
	pgRepo *postgres.Repository,
	redisRepo *redis.Repository,
	esRepo *elasticsearch.Repository,
	analyticsService *analytics.Service,
) {
	startTime := time.Now()
	log.Info().Msg("Starting DeFiLlama fetch job")

	// Fetch pools from API
	pools, err := client.FetchPools(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch pools from DeFiLlama")
		return
	}

	log.Info().Int("count", len(pools)).Msg("Fetched pools from DeFiLlama")

	// Filter pools by minimum TVL
	filteredPools := make([]defillama.Pool, 0)
	for _, p := range pools {
		if p.TVLUsd >= cfg.Worker.MinTVLThreshold {
			filteredPools = append(filteredPools, p)
		}
	}

	log.Info().
		Int("total", len(pools)).
		Int("filtered", len(filteredPools)).
		Float64("min_tvl", cfg.Worker.MinTVLThreshold).
		Msg("Filtered pools by TVL")

	// Convert to internal models and calculate scores
	modelPools := make([]models.Pool, 0, len(filteredPools))
	for _, p := range filteredPools {
		pool := defillama.ToPoolModel(p)

		// Calculate opportunity score
		pool.Score = analyticsService.CalculateScore(&pool)

		modelPools = append(modelPools, pool)
	}

	// Store in PostgreSQL (batch upsert)
	for _, pool := range modelPools {
		if err := pgRepo.UpsertPool(ctx, &pool); err != nil {
			log.Warn().Err(err).Str("pool_id", pool.ID).Msg("Failed to upsert pool")
		}

		// Record historical data point
		historical := &models.HistoricalAPY{
			PoolID:    pool.ID,
			Timestamp: time.Now().UTC(),
			APY:       pool.APY,
			TVL:       pool.TVL,
			APYBase:   pool.APYBase,
			APYReward: pool.APYReward,
		}
		if err := pgRepo.InsertHistoricalAPY(ctx, historical); err != nil {
			log.Warn().Err(err).Str("pool_id", pool.ID).Msg("Failed to insert historical APY")
		}
	}

	// Index in ElasticSearch (bulk)
	if err := esRepo.BulkIndexPools(ctx, modelPools); err != nil {
		log.Warn().Err(err).Msg("Failed to bulk index pools in ElasticSearch")
	}

	// Cache in Redis
	if err := redisRepo.SetMultiplePools(ctx, modelPools, 300); err != nil {
		log.Warn().Err(err).Msg("Failed to cache pools in Redis")
	}

	// Invalidate list caches
	if err := redisRepo.InvalidateAllPoolsCache(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to invalidate pools cache")
	}
	if err := redisRepo.InvalidateStatsCache(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to invalidate stats cache")
	}

	// Publish updates for WebSocket clients
	for _, pool := range modelPools {
		if err := redisRepo.PublishPoolUpdate(ctx, &pool); err != nil {
			log.Debug().Err(err).Str("pool_id", pool.ID).Msg("Failed to publish pool update")
		}
	}

	duration := time.Since(startTime)
	log.Info().
		Int("pools_processed", len(modelPools)).
		Dur("duration", duration).
		Msg("DeFiLlama fetch job completed")
}

// runCoinGeckoJob fetches token prices from CoinGecko
func runCoinGeckoJob(
	ctx context.Context,
	client *coingecko.Client,
	redisRepo *redis.Repository,
) {
	startTime := time.Now()
	log.Info().Msg("Starting CoinGecko fetch job")

	// Fetch prices for common tokens
	tokens := []string{
		"ethereum", "bitcoin", "tether", "usd-coin", "binance-coin",
		"matic-network", "avalanche-2", "fantom", "arbitrum", "optimism",
	}

	prices, err := client.FetchPrices(ctx, tokens)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch prices from CoinGecko")
		return
	}

	// Cache prices in Redis (15 minute TTL)
	if err := redisRepo.SetMultipleTokenPrices(ctx, prices, 900); err != nil {
		log.Warn().Err(err).Msg("Failed to cache token prices")
	}

	duration := time.Since(startTime)
	log.Info().
		Int("tokens_fetched", len(prices)).
		Dur("duration", duration).
		Msg("CoinGecko fetch job completed")
}

// runOpportunityDetectionJob analyzes pools for opportunities
func runOpportunityDetectionJob(
	ctx context.Context,
	service *opportunity.Service,
	pgRepo *postgres.Repository,
	redisRepo *redis.Repository,
) {
	startTime := time.Now()
	log.Info().Msg("Starting opportunity detection job")

	// Deactivate expired opportunities first
	if err := pgRepo.DeactivateExpiredOpportunities(ctx); err != nil {
		log.Warn().Err(err).Msg("Failed to deactivate expired opportunities")
	}

	// Detect yield gap opportunities
	yieldGaps, err := service.DetectYieldGaps(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to detect yield gaps")
	} else {
		log.Info().Int("count", len(yieldGaps)).Msg("Detected yield gap opportunities")

		// Save and publish alerts for new opportunities
		for _, opp := range yieldGaps {
			if err := pgRepo.UpsertOpportunity(ctx, &opp); err != nil {
				log.Warn().Err(err).Str("id", opp.ID).Msg("Failed to save opportunity")
			}
			if err := redisRepo.PublishOpportunityAlert(ctx, &opp); err != nil {
				log.Debug().Err(err).Msg("Failed to publish opportunity alert")
			}
		}
	}

	// Detect trending pools
	trending, err := service.DetectTrendingPools(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to detect trending pools")
	} else {
		log.Info().Int("count", len(trending)).Msg("Detected trending pools")

		// Save trending opportunities
		for _, opp := range trending {
			if err := pgRepo.UpsertOpportunity(ctx, &opp); err != nil {
				log.Warn().Err(err).Str("id", opp.ID).Msg("Failed to save trending opportunity")
			}
		}
	}

	// Detect high-score opportunities
	highScore, err := service.DetectHighScorePools(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to detect high-score pools")
	} else {
		log.Info().Int("count", len(highScore)).Msg("Detected high-score opportunities")

		// Save high-score opportunities
		for _, opp := range highScore {
			if err := pgRepo.UpsertOpportunity(ctx, &opp); err != nil {
				log.Warn().Err(err).Str("id", opp.ID).Msg("Failed to save high-score opportunity")
			}
		}
	}

	duration := time.Since(startTime)
	log.Info().
		Dur("duration", duration).
		Msg("Opportunity detection job completed")
}


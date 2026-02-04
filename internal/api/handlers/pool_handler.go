package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Request timeout for database operations
const requestTimeout = 30 * time.Second

// ListPools returns a paginated list of pools with optional filters
// @Summary List all pools
// @Description Get a paginated list of DeFi yield pools with optional filtering and sorting
// @Tags pools
// @Accept json
// @Produce json
// @Param chain query string false "Filter by blockchain (e.g., ethereum, bsc, polygon)"
// @Param protocol query string false "Filter by protocol (e.g., aave-v3, compound)"
// @Param symbol query string false "Filter by symbol (partial match)"
// @Param minApy query number false "Minimum APY percentage"
// @Param maxApy query number false "Maximum APY percentage"
// @Param minTvl query number false "Minimum TVL in USD"
// @Param maxTvl query number false "Maximum TVL in USD"
// @Param minScore query number false "Minimum risk-adjusted score (0-100)"
// @Param stablecoin query boolean false "Filter stablecoin pools only"
// @Param sortBy query string false "Sort field (apy, tvl, score)" default(tvl)
// @Param sortOrder query string false "Sort order (asc, desc)" default(desc)
// @Param limit query integer false "Number of results per page" default(50) maximum(100)
// @Param offset query integer false "Offset for pagination" default(0)
// @Success 200 {object} models.PoolListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ValidationErrors
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/pools [get]
func (h *Handler) ListPools(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), requestTimeout)
	defer cancel()

	// Parse and validate filter parameters
	filter, validationErrors := ParsePoolFilter(c)
	if len(validationErrors) > 0 {
		return SendValidationError(c, validationErrors)
	}

	// Build cache key
	cacheKey := buildPoolsCacheKey(filter)

	// Try cache first
	cached, err := h.redis.GetPoolsCache(ctx, cacheKey)
	if err == nil && cached != nil {
		log.Debug().Str("cache_key", cacheKey).Msg("Cache hit for pools")
		return c.JSON(cached)
	}

	// Fetch from ElasticSearch for fast filtering
	pools, total, err := h.es.SearchPools(ctx, filter)
	if err != nil || total == 0 {
		if err != nil {
			log.Warn().Err(err).Msg("ElasticSearch query failed, falling back to PostgreSQL")
		} else {
			log.Debug().Msg("ElasticSearch returned no results, falling back to PostgreSQL")
		}
		// Fallback to PostgreSQL
		pools, total, err = h.pg.ListPools(ctx, filter)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch pools from database")
			return SendError(c, ErrInternalServer.WithDetails("Failed to fetch pools"))
		}
	}

	response := models.PoolListResponse{
		Data:    pools,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: int64(filter.Offset+len(pools)) < total,
	}

	// Cache for 30 seconds
	if err := h.redis.SetPoolsCache(ctx, cacheKey, &response, 30); err != nil {
		log.Debug().Err(err).Msg("Failed to cache pools response")
	}

	return c.JSON(response)
}

// GetPool returns a specific pool by ID
// @Summary Get pool by ID
// @Description Get detailed information about a specific DeFi yield pool
// @Tags pools
// @Accept json
// @Produce json
// @Param id path string true "Pool ID"
// @Success 200 {object} models.Pool
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/pools/{id} [get]
func (h *Handler) GetPool(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), requestTimeout)
	defer cancel()
	poolID := c.Params("id")

	// Validate pool ID
	if errors := ValidatePoolID(poolID); len(errors) > 0 {
		return SendValidationError(c, errors)
	}

	// Try cache first
	cached, err := h.redis.GetPool(ctx, poolID)
	if err == nil && cached != nil {
		log.Debug().Str("pool_id", poolID).Msg("Cache hit for pool")
		return c.JSON(cached)
	}

	// Fetch from database
	pool, err := h.pg.GetPool(ctx, poolID)
	if err != nil {
		log.Debug().Err(err).Str("pool_id", poolID).Msg("Pool not found")
		return SendError(c, ErrNotFound.WithDetails(fmt.Sprintf("Pool '%s' not found", poolID)))
	}

	// Cache for 1 minute
	if err := h.redis.SetPool(ctx, pool, 60); err != nil {
		log.Debug().Err(err).Msg("Failed to cache pool")
	}

	return c.JSON(pool)
}

// GetPoolHistory returns historical APY data for a pool
// @Summary Get pool APY history
// @Description Get historical APY and TVL data for charting
// @Tags pools
// @Accept json
// @Produce json
// @Param id path string true "Pool ID"
// @Param period query string false "Time period (1h, 24h, 7d, 30d)" default(24h)
// @Success 200 {object} models.PoolHistoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/pools/{id}/history [get]
func (h *Handler) GetPoolHistory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), requestTimeout)
	defer cancel()
	poolID := c.Params("id")
	period := c.Query("period", "24h")

	// Validate pool ID
	if errors := ValidatePoolID(poolID); len(errors) > 0 {
		return SendValidationError(c, errors)
	}

	// Validate period
	if errors := ValidatePeriod(period); len(errors) > 0 {
		return SendValidationError(c, errors)
	}

	// Fetch historical data from TimescaleDB
	history, err := h.pg.GetPoolHistory(ctx, poolID, period)
	if err != nil {
		log.Error().Err(err).
			Str("pool_id", poolID).
			Str("period", period).
			Msg("Failed to fetch pool history")
		return SendError(c, ErrInternalServer.WithDetails("Failed to fetch pool history"))
	}

	response := models.PoolHistoryResponse{
		PoolID:     poolID,
		Period:     period,
		DataPoints: history,
	}

	return c.JSON(response)
}

// buildPoolsCacheKey creates a cache key from filter parameters
func buildPoolsCacheKey(filter models.PoolFilter) string {
	stablecoin := ""
	if filter.StableCoin != nil {
		if *filter.StableCoin {
			stablecoin = "true"
		} else {
			stablecoin = "false"
		}
	}
	return fmt.Sprintf("pools:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%d:%d",
		filter.Chain,
		filter.Protocol,
		filter.Symbol,
		filter.Search,
		filter.MinAPY.String(),
		filter.MaxAPY.String(),
		filter.MinTVL.String(),
		filter.MaxTVL.String(),
		filter.MinScore.String(),
		stablecoin,
		filter.SortBy,
		filter.SortOrder,
		filter.Limit,
		filter.Offset,
	)
}

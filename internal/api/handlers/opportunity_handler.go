package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// ListOpportunities returns detected yield farming opportunities
// @Summary List opportunities
// @Description Get a list of detected yield farming opportunities with filtering
// @Tags opportunities
// @Accept json
// @Produce json
// @Param type query string false "Opportunity type (yield-gap, trending, high-score)"
// @Param riskLevel query string false "Risk level (low, medium, high)"
// @Param chain query string false "Filter by blockchain"
// @Param asset query string false "Filter by asset (e.g., USDC, ETH)"
// @Param minProfit query number false "Minimum potential profit percentage"
// @Param minScore query number false "Minimum opportunity score"
// @Param activeOnly query boolean false "Show only active opportunities" default(true)
// @Param sortBy query string false "Sort field (score, profit, apy, detected_at)" default(score)
// @Param sortOrder query string false "Sort order (asc, desc)" default(desc)
// @Param limit query integer false "Number of results per page" default(50) maximum(100)
// @Param offset query integer false "Offset for pagination" default(0)
// @Success 200 {object} models.OpportunityListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ValidationErrors
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/opportunities [get]
func (h *Handler) ListOpportunities(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Parse and validate filter parameters
	filter, validationErrors := ParseOpportunityFilter(c)
	if len(validationErrors) > 0 {
		return SendValidationError(c, validationErrors)
	}

	// Build cache key
	cacheKey := buildOpportunitiesCacheKey(filter)

	// Try cache first
	cached, err := h.redis.GetOpportunitiesCache(ctx, cacheKey)
	if err == nil && cached != nil {
		log.Debug().Str("cache_key", cacheKey).Msg("Cache hit for opportunities")
		return c.JSON(cached)
	}

	// Fetch from database
	opportunities, total, err := h.pg.ListOpportunities(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch opportunities")
		return SendError(c, ErrInternalServer.WithDetails("Failed to fetch opportunities"))
	}

	response := models.OpportunityListResponse{
		Data:    opportunities,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: int64(filter.Offset+len(opportunities)) < total,
	}

	// Cache for 1 minute
	if err := h.redis.SetOpportunitiesCache(ctx, cacheKey, &response, 60); err != nil {
		log.Debug().Err(err).Msg("Failed to cache opportunities response")
	}

	return c.JSON(response)
}

// GetTrendingPools returns pools with significantly increasing APY
// @Summary Get trending pools
// @Description Get pools with rapidly increasing APY in the last 24 hours
// @Tags opportunities
// @Accept json
// @Produce json
// @Param chain query string false "Filter by blockchain"
// @Param minGrowth query number false "Minimum APY growth percentage" default(10)
// @Param limit query integer false "Number of results" default(20) maximum(50)
// @Param offset query integer false "Offset for pagination" default(0)
// @Success 200 {object} TrendingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/opportunities/trending [get]
func (h *Handler) GetTrendingPools(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	chain := c.Query("chain")
	minGrowthStr := c.Query("minGrowth", "10")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Validate parameters
	var validationErrors []ValidationError

	minGrowth, err := decimal.NewFromString(minGrowthStr)
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "minGrowth",
			Message: "must be a valid number",
		})
	}

	if limit > 50 {
		limit = 50
	}
	if limit < 1 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	if len(validationErrors) > 0 {
		return SendValidationError(c, validationErrors)
	}

	// Try cache first
	cacheKey := fmt.Sprintf("trending:%s:%.1f", chain, minGrowth.InexactFloat64())
	cached, err := h.redis.GetTrendingCache(ctx, cacheKey)
	if err == nil && cached != nil {
		log.Debug().Str("cache_key", cacheKey).Msg("Cache hit for trending pools")
		return c.JSON(TrendingResponse{
			Data:   cached,
			Limit:  limit,
			Offset: offset,
		})
	}

	// Fetch trending pools
	trending, err := h.pg.GetTrendingPools(ctx, chain, minGrowth, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch trending pools")
		return SendError(c, ErrInternalServer.WithDetails("Failed to fetch trending pools"))
	}

	// Cache for 2 minutes
	if err := h.redis.SetTrendingCache(ctx, cacheKey, trending, 120); err != nil {
		log.Debug().Err(err).Msg("Failed to cache trending pools")
	}

	return c.JSON(TrendingResponse{
		Data:   trending,
		Limit:  limit,
		Offset: offset,
	})
}

// TrendingResponse is the response for trending pools endpoint
type TrendingResponse struct {
	Data   []models.TrendingPool `json:"data"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

// buildOpportunitiesCacheKey creates a cache key for opportunities
func buildOpportunitiesCacheKey(filter models.OpportunityFilter) string {
	return fmt.Sprintf("opportunities:%s:%s:%s:%s:%s:%s:%s:%t:%d:%d",
		filter.Type,
		filter.RiskLevel,
		filter.Chain,
		filter.Asset,
		filter.MinProfit.String(),
		filter.SortBy,
		filter.SortOrder,
		filter.ActiveOnly,
		filter.Limit,
		filter.Offset,
	)
}

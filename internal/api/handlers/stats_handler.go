package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// ListChains returns all supported blockchain networks with statistics
// GET /api/v1/chains
func (h *Handler) ListChains(c *fiber.Ctx) error {
	ctx := c.Context()

	// Try cache first
	cached, err := h.redis.GetChainsCache(ctx)
	if err == nil && cached != nil {
		return c.JSON(cached)
	}

	// Fetch from database
	chains, err := h.pg.ListChains(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch chains")
	}

	response := models.ChainListResponse{
		Data:  chains,
		Total: len(chains),
	}

	// Cache for 5 minutes (chain data doesn't change often)
	_ = h.redis.SetChainsCache(ctx, &response, 300)

	return c.JSON(response)
}

// ListProtocols returns all DeFi protocols with statistics
// GET /api/v1/protocols
// Query params: chain, category, sortBy, sortOrder, limit, offset
func (h *Handler) ListProtocols(c *fiber.Ctx) error {
	ctx := c.Context()

	filter := models.ProtocolFilter{
		Chain:     c.Query("chain"),
		Category:  c.Query("category"),
		SortBy:    c.Query("sortBy", "tvl"),
		SortOrder: c.Query("sortOrder", "desc"),
		Limit:     c.QueryInt("limit", 50),
		Offset:    c.QueryInt("offset", 0),
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// Try cache first
	cacheKey := "protocols:" + filter.Chain + ":" + filter.Category
	cached, err := h.redis.GetProtocolsCache(ctx, cacheKey)
	if err == nil && cached != nil {
		return c.JSON(cached)
	}

	// Fetch from database
	protocols, total, err := h.pg.ListProtocols(ctx, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch protocols")
	}

	response := models.ProtocolListResponse{
		Data:    protocols,
		Total:   total,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
		HasMore: int64(filter.Offset+len(protocols)) < total,
	}

	// Cache for 5 minutes
	_ = h.redis.SetProtocolsCache(ctx, cacheKey, &response, 300)

	return c.JSON(response)
}

// GetStats returns overall platform statistics
// GET /api/v1/stats
func (h *Handler) GetStats(c *fiber.Ctx) error {
	ctx := c.Context()

	// Try cache first
	cached, err := h.redis.GetStatsCache(ctx)
	if err == nil && cached != nil {
		return c.JSON(cached)
	}

	// Fetch fresh stats from database
	stats, err := h.pg.GetPlatformStats(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch statistics")
	}

	// Cache for 2 minutes (stats should be relatively fresh)
	_ = h.redis.SetStatsCache(ctx, stats, 120)

	return c.JSON(stats)
}

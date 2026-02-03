// Package graphql provides GraphQL API implementation for the DeFi Yield Aggregator.
// It offers flexible queries for complex data needs as an alternative to REST.
package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/elasticsearch"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/postgres"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
)

// Resolver handles GraphQL query resolution
type Resolver struct {
	pg        *postgres.Repository
	redis     *redis.Repository
	es        *elasticsearch.Repository
	startTime time.Time
}

// NewResolver creates a new GraphQL resolver
func NewResolver(pg *postgres.Repository, redis *redis.Repository, es *elasticsearch.Repository) *Resolver {
	return &Resolver{
		pg:        pg,
		redis:     redis,
		es:        es,
		startTime: time.Now(),
	}
}

// GraphQL request/response types
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{}      `json:"data,omitempty"`
	Errors []GraphQLError   `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message   string        `json:"message"`
	Locations []Location    `json:"locations,omitempty"`
	Path      []interface{} `json:"path,omitempty"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Handle processes GraphQL requests
// This is a simplified implementation - for production, use gqlgen or graphql-go
func (r *Resolver) Handle(c *fiber.Ctx) error {
	var req GraphQLRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(GraphQLResponse{
			Errors: []GraphQLError{{Message: "Invalid request body"}},
		})
	}

	ctx := c.Context()
	data, errors := r.executeQuery(ctx, req)

	response := GraphQLResponse{
		Data:   data,
		Errors: errors,
	}

	return c.JSON(response)
}

// executeQuery parses and executes GraphQL queries
// This is a simplified query executor - supports common operations
func (r *Resolver) executeQuery(ctx context.Context, req GraphQLRequest) (interface{}, []GraphQLError) {
	// For a full implementation, use gqlgen or graphql-go library
	// This simplified version handles common query patterns

	data := make(map[string]interface{})
	var errors []GraphQLError

	// Parse the query to determine what's being requested
	// This is a simplified parser - production should use proper GraphQL parsing

	if containsQuery(req.Query, "pools") && !containsQuery(req.Query, "trendingPools") {
		pools, err := r.resolvePools(ctx, req.Variables)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["pools"] = pools
		}
	}

	if containsQuery(req.Query, "pool(") {
		pool, err := r.resolvePool(ctx, req.Variables)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["pool"] = pool
		}
	}

	if containsQuery(req.Query, "opportunities") && !containsQuery(req.Query, "activeOpportunities") {
		opps, err := r.resolveOpportunities(ctx, req.Variables)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["opportunities"] = opps
		}
	}

	if containsQuery(req.Query, "trendingPools") {
		trending, err := r.resolveTrendingPools(ctx, req.Variables)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["trendingPools"] = trending
		}
	}

	if containsQuery(req.Query, "chains") {
		chains, err := r.resolveChains(ctx)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["chains"] = chains
		}
	}

	if containsQuery(req.Query, "protocols") {
		protocols, err := r.resolveProtocols(ctx, req.Variables)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["protocols"] = protocols
		}
	}

	if containsQuery(req.Query, "stats") {
		stats, err := r.resolveStats(ctx)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["stats"] = stats
		}
	}

	if containsQuery(req.Query, "health") {
		health, err := r.resolveHealth(ctx)
		if err != nil {
			errors = append(errors, GraphQLError{Message: err.Error()})
		} else {
			data["health"] = health
		}
	}

	return data, errors
}

// Pool resolvers

func (r *Resolver) resolvePools(ctx context.Context, vars map[string]interface{}) (interface{}, error) {
	filter := parsePoolFilterFromVars(vars)

	pools, total, err := r.pg.ListPools(ctx, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]map[string]interface{}, len(pools))
	for i, pool := range pools {
		edges[i] = map[string]interface{}{
			"node":   poolToGraphQL(pool),
			"cursor": encodeCursor(filter.Offset + i),
		}
	}

	return map[string]interface{}{
		"edges": edges,
		"pageInfo": map[string]interface{}{
			"hasNextPage":     int64(filter.Offset+len(pools)) < total,
			"hasPreviousPage": filter.Offset > 0,
			"startCursor":     encodeCursor(filter.Offset),
			"endCursor":       encodeCursor(filter.Offset + len(pools) - 1),
		},
		"totalCount": total,
	}, nil
}

func (r *Resolver) resolvePool(ctx context.Context, vars map[string]interface{}) (interface{}, error) {
	id, ok := vars["id"].(string)
	if !ok {
		return nil, fmt.Errorf("pool id is required")
	}

	pool, err := r.pg.GetPool(ctx, id)
	if err != nil {
		return nil, err
	}

	return poolToGraphQL(*pool), nil
}

// Opportunity resolvers

func (r *Resolver) resolveOpportunities(ctx context.Context, vars map[string]interface{}) (interface{}, error) {
	filter := parseOpportunityFilterFromVars(vars)

	opps, total, err := r.pg.ListOpportunities(ctx, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]map[string]interface{}, len(opps))
	for i, opp := range opps {
		edges[i] = map[string]interface{}{
			"node":   opportunityToGraphQL(opp),
			"cursor": encodeCursor(filter.Offset + i),
		}
	}

	return map[string]interface{}{
		"edges": edges,
		"pageInfo": map[string]interface{}{
			"hasNextPage":     int64(filter.Offset+len(opps)) < total,
			"hasPreviousPage": filter.Offset > 0,
		},
		"totalCount": total,
	}, nil
}

func (r *Resolver) resolveTrendingPools(ctx context.Context, vars map[string]interface{}) (interface{}, error) {
	chain := ""
	if c, ok := vars["chain"].(string); ok {
		chain = c
	}

	minGrowth := decimal.NewFromFloat(10)
	if mg, ok := vars["minGrowth"].(float64); ok {
		minGrowth = decimal.NewFromFloat(mg)
	}

	limit := 20
	if l, ok := vars["limit"].(float64); ok {
		limit = int(l)
	}

	trending, err := r.pg.GetTrendingPools(ctx, chain, minGrowth, limit, 0)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(trending))
	for i, tp := range trending {
		result[i] = map[string]interface{}{
			"pool":         poolToGraphQL(*tp.Pool),
			"apyGrowth1h":  tp.APYGrowth1H.String(),
			"apyGrowth24h": tp.APYGrowth24H.String(),
			"apyGrowth7d":  tp.APYGrowth7D.String(),
			"trendScore":   tp.TrendScore.String(),
		}
	}

	return result, nil
}

// Stats resolvers

func (r *Resolver) resolveChains(ctx context.Context) (interface{}, error) {
	chains, err := r.pg.ListChains(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(chains))
	for i, c := range chains {
		result[i] = map[string]interface{}{
			"name":        c.Name,
			"displayName": c.DisplayName,
			"poolCount":   c.PoolCount,
			"totalTvl":    c.TotalTVL.String(),
			"averageApy":  c.AverageAPY.String(),
			"maxApy":      c.MaxAPY.String(),
		}
	}

	return result, nil
}

func (r *Resolver) resolveProtocols(ctx context.Context, vars map[string]interface{}) (interface{}, error) {
	filter := models.ProtocolFilter{
		Limit:  50,
		Offset: 0,
	}

	if chain, ok := vars["chain"].(string); ok {
		filter.Chain = chain
	}

	protocols, total, err := r.pg.ListProtocols(ctx, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]map[string]interface{}, len(protocols))
	for i, p := range protocols {
		edges[i] = map[string]interface{}{
			"node": map[string]interface{}{
				"name":        p.Name,
				"displayName": p.DisplayName,
				"chains":      p.Chains,
				"poolCount":   p.PoolCount,
				"totalTvl":    p.TotalTVL.String(),
				"averageApy":  p.AverageAPY.String(),
				"maxApy":      p.MaxAPY.String(),
			},
			"cursor": encodeCursor(i),
		}
	}

	return map[string]interface{}{
		"edges":      edges,
		"totalCount": total,
	}, nil
}

func (r *Resolver) resolveStats(ctx context.Context) (interface{}, error) {
	stats, err := r.pg.GetPlatformStats(ctx)
	if err != nil {
		return nil, err
	}

	// Convert tvlByChain
	tvlByChain := make([]map[string]interface{}, 0)
	for chain, tvl := range stats.TVLByChain {
		tvlByChain = append(tvlByChain, map[string]interface{}{
			"chain": chain,
			"tvl":   tvl.String(),
		})
	}

	// Convert poolsByChain
	poolsByChain := make([]map[string]interface{}, 0)
	for chain, count := range stats.PoolsByChain {
		poolsByChain = append(poolsByChain, map[string]interface{}{
			"chain": chain,
			"count": count,
		})
	}

	return map[string]interface{}{
		"totalPools":          stats.TotalPools,
		"totalTvl":            stats.TotalTVL.String(),
		"averageApy":          stats.AverageAPY.String(),
		"medianApy":           stats.MedianAPY.String(),
		"maxApy":              stats.MaxAPY.String(),
		"totalChains":         stats.TotalChains,
		"totalProtocols":      stats.TotalProtocols,
		"activeOpportunities": stats.ActiveOpportunities,
		"lastUpdated":         stats.LastUpdated,
		"tvlByChain":          tvlByChain,
		"poolsByChain":        poolsByChain,
		"apyDistribution": map[string]interface{}{
			"range0to1":    stats.APYDistribution.Range0to1,
			"range1to5":    stats.APYDistribution.Range1to5,
			"range5to10":   stats.APYDistribution.Range5to10,
			"range10to25":  stats.APYDistribution.Range10to25,
			"range25to50":  stats.APYDistribution.Range25to50,
			"range50to100": stats.APYDistribution.Range50to100,
			"range100plus": stats.APYDistribution.Range100Plus,
		},
	}, nil
}

func (r *Resolver) resolveHealth(ctx context.Context) (interface{}, error) {
	health := map[string]interface{}{
		"status":    "HEALTHY",
		"version":   "1.0.0",
		"uptime":    time.Since(r.startTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"services":  make(map[string]interface{}),
	}

	services := health["services"].(map[string]interface{})

	// Check PostgreSQL
	pgStart := time.Now()
	pgErr := r.pg.Ping(ctx)
	services["postgresql"] = map[string]interface{}{
		"status":  boolToStatus(pgErr == nil),
		"latency": time.Since(pgStart).String(),
		"message": errToMessage(pgErr),
	}

	// Check Redis
	redisStart := time.Now()
	redisErr := r.redis.Ping(ctx)
	services["redis"] = map[string]interface{}{
		"status":  boolToStatus(redisErr == nil),
		"latency": time.Since(redisStart).String(),
		"message": errToMessage(redisErr),
	}

	// Check ElasticSearch
	esStart := time.Now()
	esErr := r.es.Ping(ctx)
	services["elasticsearch"] = map[string]interface{}{
		"status":  boolToStatus(esErr == nil),
		"latency": time.Since(esStart).String(),
		"message": errToMessage(esErr),
	}

	// Determine overall status
	if pgErr != nil || redisErr != nil {
		health["status"] = "UNHEALTHY"
	} else if esErr != nil {
		health["status"] = "DEGRADED"
	}

	return health, nil
}

// Helper functions

func containsQuery(query, field string) bool {
	return len(query) > 0 && (contains(query, field) || contains(query, "{"+field))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func parsePoolFilterFromVars(vars map[string]interface{}) models.PoolFilter {
	filter := models.PoolFilter{
		Limit:     50,
		Offset:    0,
		SortBy:    "tvl",
		SortOrder: "desc",
	}

	if filterVar, ok := vars["filter"].(map[string]interface{}); ok {
		if chain, ok := filterVar["chain"].(string); ok {
			filter.Chain = chain
		}
		if protocol, ok := filterVar["protocol"].(string); ok {
			filter.Protocol = protocol
		}
		if symbol, ok := filterVar["symbol"].(string); ok {
			filter.Symbol = symbol
		}
		if minApy, ok := filterVar["minApy"].(float64); ok {
			filter.MinAPY = decimal.NewFromFloat(minApy)
		}
		if maxApy, ok := filterVar["maxApy"].(float64); ok {
			filter.MaxAPY = decimal.NewFromFloat(maxApy)
		}
		if minTvl, ok := filterVar["minTvl"].(float64); ok {
			filter.MinTVL = decimal.NewFromFloat(minTvl)
		}
		if stablecoin, ok := filterVar["stablecoin"].(bool); ok {
			filter.StableCoin = &stablecoin
		}
	}

	if paginationVar, ok := vars["pagination"].(map[string]interface{}); ok {
		if limit, ok := paginationVar["limit"].(float64); ok {
			filter.Limit = int(limit)
		}
		if offset, ok := paginationVar["offset"].(float64); ok {
			filter.Offset = int(offset)
		}
	}

	return filter
}

func parseOpportunityFilterFromVars(vars map[string]interface{}) models.OpportunityFilter {
	filter := models.OpportunityFilter{
		Limit:      50,
		Offset:     0,
		ActiveOnly: true,
		SortBy:     "score",
		SortOrder:  "desc",
	}

	if filterVar, ok := vars["filter"].(map[string]interface{}); ok {
		if t, ok := filterVar["type"].(string); ok {
			filter.Type = models.OpportunityType(t)
		}
		if risk, ok := filterVar["riskLevel"].(string); ok {
			filter.RiskLevel = models.RiskLevel(risk)
		}
		if chain, ok := filterVar["chain"].(string); ok {
			filter.Chain = chain
		}
	}

	return filter
}

func poolToGraphQL(pool models.Pool) map[string]interface{} {
	return map[string]interface{}{
		"id":               pool.ID,
		"chain":            pool.Chain,
		"protocol":         pool.Protocol,
		"symbol":           pool.Symbol,
		"tvl":              pool.TVL.String(),
		"apy":              pool.APY.String(),
		"apyBase":          pool.APYBase.String(),
		"apyReward":        pool.APYReward.String(),
		"rewardTokens":     pool.RewardTokens,
		"underlyingTokens": pool.UnderlyingTokens,
		"poolMeta":         pool.PoolMeta,
		"il7d":             pool.IL7D.String(),
		"apyMean30d":       pool.APYMean30D.String(),
		"volumeUsd1d":      pool.VolumeUSD1D.String(),
		"volumeUsd7d":      pool.VolumeUSD7D.String(),
		"score":            pool.Score.String(),
		"apyChange1h":      pool.APYChange1H.String(),
		"apyChange24h":     pool.APYChange24H.String(),
		"apyChange7d":      pool.APYChange7D.String(),
		"stablecoin":       pool.StableCoin,
		"exposure":         pool.Exposure,
		"createdAt":        pool.CreatedAt.Format(time.RFC3339),
		"updatedAt":        pool.UpdatedAt.Format(time.RFC3339),
	}
}

func opportunityToGraphQL(opp models.Opportunity) map[string]interface{} {
	result := map[string]interface{}{
		"id":              opp.ID,
		"type":            string(opp.Type),
		"title":           opp.Title,
		"description":     opp.Description,
		"asset":           opp.Asset,
		"chain":           opp.Chain,
		"apyDifference":   opp.APYDifference.String(),
		"apyGrowth":       opp.APYGrowth.String(),
		"currentApy":      opp.CurrentAPY.String(),
		"potentialProfit": opp.PotentialProfit.String(),
		"tvl":             opp.TVL.String(),
		"riskLevel":       string(opp.RiskLevel),
		"score":           opp.Score.String(),
		"isActive":        opp.IsActive,
		"detectedAt":      opp.DetectedAt.Format(time.RFC3339),
		"lastSeenAt":      opp.LastSeenAt.Format(time.RFC3339),
		"expiresAt":       opp.ExpiresAt.Format(time.RFC3339),
		"createdAt":       opp.CreatedAt.Format(time.RFC3339),
		"updatedAt":       opp.UpdatedAt.Format(time.RFC3339),
	}

	return result
}

func encodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

func boolToStatus(ok bool) string {
	if ok {
		return "UP"
	}
	return "DOWN"
}

func errToMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// Ensure json is used
var _ = json.Marshal

// Ensure log is used
var _ = log.Debug

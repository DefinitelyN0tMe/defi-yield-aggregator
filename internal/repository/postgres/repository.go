// Package postgres provides PostgreSQL database operations for the DeFi Yield Aggregator.
// It uses pgx for high-performance PostgreSQL connectivity and supports TimescaleDB
// for efficient time-series data storage.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Repository handles all PostgreSQL database operations
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new PostgreSQL repository with connection pooling
func NewRepository(ctx context.Context, cfg config.PostgresConfig) (*Repository, error) {
	// Build connection string
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
		cfg.MaxConnections, cfg.MaxIdleConnections, cfg.ConnectionMaxLifetime,
	)

	// Configure the connection pool
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Add logging for query execution (development only)
	poolConfig.ConnConfig.Tracer = &queryTracer{}

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{pool: pool}, nil
}

// Close closes the database connection pool
func (r *Repository) Close() {
	r.pool.Close()
}

// Ping checks if the database connection is alive
func (r *Repository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

// queryTracer implements pgx.QueryTracer for logging queries
type queryTracer struct{}

func (t *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Debug().
		Str("sql", data.SQL).
		Interface("args", data.Args).
		Msg("Executing query")
	return ctx
}

func (t *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		log.Error().Err(data.Err).Msg("Query failed")
	}
}

// =============================================================================
// Pool Operations
// =============================================================================

// ListPools returns a paginated list of pools with optional filters
func (r *Repository) ListPools(ctx context.Context, filter models.PoolFilter) ([]models.Pool, int64, error) {
	// Build dynamic query based on filters
	query := `
		SELECT
			id, chain, protocol, symbol, tvl, apy, apy_base, apy_reward,
			reward_tokens, underlying_tokens, pool_meta, il_7d, apy_mean_30d,
			volume_usd_1d, volume_usd_7d, score, apy_change_1h, apy_change_24h,
			apy_change_7d, stablecoin, exposure, created_at, updated_at
		FROM pools
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM pools WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	// Apply filters (using ILIKE for case-insensitive matching)
	if filter.Chain != "" {
		argCount++
		query += fmt.Sprintf(" AND LOWER(chain) = LOWER($%d)", argCount)
		countQuery += fmt.Sprintf(" AND LOWER(chain) = LOWER($%d)", argCount)
		args = append(args, filter.Chain)
	}

	if filter.Protocol != "" {
		argCount++
		query += fmt.Sprintf(" AND LOWER(protocol) = LOWER($%d)", argCount)
		countQuery += fmt.Sprintf(" AND LOWER(protocol) = LOWER($%d)", argCount)
		args = append(args, filter.Protocol)
	}

	if filter.Symbol != "" {
		argCount++
		query += fmt.Sprintf(" AND symbol ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND symbol ILIKE $%d", argCount)
		args = append(args, "%"+filter.Symbol+"%")
	}

	// Search across multiple fields (symbol, protocol, chain, pool_meta)
	if filter.Search != "" {
		argCount++
		searchPattern := "%" + filter.Search + "%"
		query += fmt.Sprintf(" AND (symbol ILIKE $%d OR protocol ILIKE $%d OR chain ILIKE $%d OR pool_meta ILIKE $%d)", argCount, argCount, argCount, argCount)
		countQuery += fmt.Sprintf(" AND (symbol ILIKE $%d OR protocol ILIKE $%d OR chain ILIKE $%d OR pool_meta ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, searchPattern)
	}

	if !filter.MinAPY.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND apy >= $%d", argCount)
		countQuery += fmt.Sprintf(" AND apy >= $%d", argCount)
		args = append(args, filter.MinAPY)
	}

	if !filter.MaxAPY.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND apy <= $%d", argCount)
		countQuery += fmt.Sprintf(" AND apy <= $%d", argCount)
		args = append(args, filter.MaxAPY)
	}

	if !filter.MinTVL.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND tvl >= $%d", argCount)
		countQuery += fmt.Sprintf(" AND tvl >= $%d", argCount)
		args = append(args, filter.MinTVL)
	}

	if !filter.MaxTVL.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND tvl <= $%d", argCount)
		countQuery += fmt.Sprintf(" AND tvl <= $%d", argCount)
		args = append(args, filter.MaxTVL)
	}

	if filter.StableCoin != nil {
		argCount++
		query += fmt.Sprintf(" AND stablecoin = $%d", argCount)
		countQuery += fmt.Sprintf(" AND stablecoin = $%d", argCount)
		args = append(args, *filter.StableCoin)
	}

	// Get total count
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pools: %w", err)
	}

	// Add sorting
	sortColumn := "tvl"
	switch filter.SortBy {
	case "apy":
		sortColumn = "apy"
	case "score":
		sortColumn = "score"
	case "tvl":
		sortColumn = "tvl"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filter.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filter.Offset)

	// Execute query
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pools: %w", err)
	}
	defer rows.Close()

	pools := make([]models.Pool, 0)
	for rows.Next() {
		var pool models.Pool
		err := rows.Scan(
			&pool.ID, &pool.Chain, &pool.Protocol, &pool.Symbol,
			&pool.TVL, &pool.APY, &pool.APYBase, &pool.APYReward,
			&pool.RewardTokens, &pool.UnderlyingTokens, &pool.PoolMeta,
			&pool.IL7D, &pool.APYMean30D, &pool.VolumeUSD1D, &pool.VolumeUSD7D,
			&pool.Score, &pool.APYChange1H, &pool.APYChange24H, &pool.APYChange7D,
			&pool.StableCoin, &pool.Exposure, &pool.CreatedAt, &pool.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pool: %w", err)
		}
		pools = append(pools, pool)
	}

	return pools, total, nil
}

// GetPool returns a single pool by ID
func (r *Repository) GetPool(ctx context.Context, id string) (*models.Pool, error) {
	query := `
		SELECT
			id, chain, protocol, symbol, tvl, apy, apy_base, apy_reward,
			reward_tokens, underlying_tokens, pool_meta, il_7d, apy_mean_30d,
			volume_usd_1d, volume_usd_7d, score, apy_change_1h, apy_change_24h,
			apy_change_7d, stablecoin, exposure, created_at, updated_at
		FROM pools
		WHERE id = $1
	`

	var pool models.Pool
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pool.ID, &pool.Chain, &pool.Protocol, &pool.Symbol,
		&pool.TVL, &pool.APY, &pool.APYBase, &pool.APYReward,
		&pool.RewardTokens, &pool.UnderlyingTokens, &pool.PoolMeta,
		&pool.IL7D, &pool.APYMean30D, &pool.VolumeUSD1D, &pool.VolumeUSD7D,
		&pool.Score, &pool.APYChange1H, &pool.APYChange24H, &pool.APYChange7D,
		&pool.StableCoin, &pool.Exposure, &pool.CreatedAt, &pool.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("pool not found")
		}
		return nil, fmt.Errorf("failed to get pool: %w", err)
	}

	return &pool, nil
}

// GetPoolHistory returns historical APY data for a pool
func (r *Repository) GetPoolHistory(ctx context.Context, poolID string, period string) ([]models.HistoricalAPY, error) {
	// Calculate time range based on period
	var interval string
	var bucketInterval string

	switch period {
	case "1h":
		interval = "1 hour"
		bucketInterval = "1 minute"
	case "24h":
		interval = "24 hours"
		bucketInterval = "5 minutes"
	case "7d":
		interval = "7 days"
		bucketInterval = "1 hour"
	case "30d":
		interval = "30 days"
		bucketInterval = "6 hours"
	default:
		interval = "24 hours"
		bucketInterval = "5 minutes"
	}

	// Use TimescaleDB time_bucket for efficient aggregation
	query := fmt.Sprintf(`
		SELECT
			pool_id,
			time_bucket('%s', timestamp) AS bucket,
			AVG(apy) AS apy,
			AVG(tvl) AS tvl,
			AVG(apy_base) AS apy_base,
			AVG(apy_reward) AS apy_reward
		FROM historical_apy
		WHERE pool_id = $1
		  AND timestamp > NOW() - INTERVAL '%s'
		GROUP BY pool_id, bucket
		ORDER BY bucket ASC
	`, bucketInterval, interval)

	rows, err := r.pool.Query(ctx, query, poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pool history: %w", err)
	}
	defer rows.Close()

	history := make([]models.HistoricalAPY, 0)
	for rows.Next() {
		var h models.HistoricalAPY
		err := rows.Scan(&h.PoolID, &h.Timestamp, &h.APY, &h.TVL, &h.APYBase, &h.APYReward)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

// UpsertPool inserts or updates a pool
func (r *Repository) UpsertPool(ctx context.Context, pool *models.Pool) error {
	query := `
		INSERT INTO pools (
			id, chain, protocol, symbol, tvl, apy, apy_base, apy_reward,
			reward_tokens, underlying_tokens, pool_meta, il_7d, apy_mean_30d,
			volume_usd_1d, volume_usd_7d, score, apy_change_1h, apy_change_24h,
			apy_change_7d, stablecoin, exposure, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)
		ON CONFLICT (id) DO UPDATE SET
			tvl = EXCLUDED.tvl,
			apy = EXCLUDED.apy,
			apy_base = EXCLUDED.apy_base,
			apy_reward = EXCLUDED.apy_reward,
			reward_tokens = EXCLUDED.reward_tokens,
			il_7d = EXCLUDED.il_7d,
			apy_mean_30d = EXCLUDED.apy_mean_30d,
			volume_usd_1d = EXCLUDED.volume_usd_1d,
			volume_usd_7d = EXCLUDED.volume_usd_7d,
			score = EXCLUDED.score,
			apy_change_1h = EXCLUDED.apy_change_1h,
			apy_change_24h = EXCLUDED.apy_change_24h,
			apy_change_7d = EXCLUDED.apy_change_7d,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		pool.ID, pool.Chain, pool.Protocol, pool.Symbol,
		pool.TVL, pool.APY, pool.APYBase, pool.APYReward,
		pool.RewardTokens, pool.UnderlyingTokens, pool.PoolMeta,
		pool.IL7D, pool.APYMean30D, pool.VolumeUSD1D, pool.VolumeUSD7D,
		pool.Score, pool.APYChange1H, pool.APYChange24H, pool.APYChange7D,
		pool.StableCoin, pool.Exposure, pool.CreatedAt, pool.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert pool: %w", err)
	}

	return nil
}

// InsertHistoricalAPY records a historical APY data point
func (r *Repository) InsertHistoricalAPY(ctx context.Context, h *models.HistoricalAPY) error {
	query := `
		INSERT INTO historical_apy (pool_id, timestamp, apy, tvl, apy_base, apy_reward)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		h.PoolID, h.Timestamp, h.APY, h.TVL, h.APYBase, h.APYReward,
	)

	if err != nil {
		return fmt.Errorf("failed to insert historical APY: %w", err)
	}

	return nil
}

// =============================================================================
// Opportunity Operations
// =============================================================================

// ListOpportunities returns opportunities based on filters
func (r *Repository) ListOpportunities(ctx context.Context, filter models.OpportunityFilter) ([]models.Opportunity, int64, error) {
	query := `
		SELECT
			id, type, title, description, source_pool_id, target_pool_id,
			pool_id, asset, chain, apy_difference, apy_growth, current_apy,
			potential_profit, tvl, risk_level, score, is_active,
			detected_at, last_seen_at, expires_at, created_at, updated_at
		FROM opportunities
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM opportunities WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	if filter.ActiveOnly {
		query += " AND is_active = true"
		countQuery += " AND is_active = true"
	}

	if filter.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		countQuery += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, filter.Type)
	}

	if filter.RiskLevel != "" {
		argCount++
		query += fmt.Sprintf(" AND risk_level = $%d", argCount)
		countQuery += fmt.Sprintf(" AND risk_level = $%d", argCount)
		args = append(args, filter.RiskLevel)
	}

	if filter.Chain != "" {
		argCount++
		query += fmt.Sprintf(" AND chain = $%d", argCount)
		countQuery += fmt.Sprintf(" AND chain = $%d", argCount)
		args = append(args, filter.Chain)
	}

	if !filter.MinProfit.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND potential_profit >= $%d", argCount)
		countQuery += fmt.Sprintf(" AND potential_profit >= $%d", argCount)
		args = append(args, filter.MinProfit)
	}

	// Get total count
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count opportunities: %w", err)
	}

	// Add sorting
	sortColumn := "score"
	switch filter.SortBy {
	case "profit":
		sortColumn = "potential_profit"
	case "apy":
		sortColumn = "current_apy"
	case "detectedAt":
		sortColumn = "detected_at"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filter.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filter.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query opportunities: %w", err)
	}
	defer rows.Close()

	opportunities := make([]models.Opportunity, 0)
	for rows.Next() {
		var o models.Opportunity
		err := rows.Scan(
			&o.ID, &o.Type, &o.Title, &o.Description,
			&o.SourcePoolID, &o.TargetPoolID, &o.PoolID,
			&o.Asset, &o.Chain, &o.APYDifference, &o.APYGrowth,
			&o.CurrentAPY, &o.PotentialProfit, &o.TVL, &o.RiskLevel,
			&o.Score, &o.IsActive, &o.DetectedAt, &o.LastSeenAt,
			&o.ExpiresAt, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan opportunity: %w", err)
		}
		opportunities = append(opportunities, o)
	}

	return opportunities, total, nil
}

// GetTrendingPools returns pools with significant APY growth
func (r *Repository) GetTrendingPools(ctx context.Context, chain string, minGrowth decimal.Decimal, limit, offset int) ([]models.TrendingPool, error) {
	query := `
		SELECT
			p.id, p.chain, p.protocol, p.symbol, p.tvl, p.apy,
			p.apy_base, p.apy_reward, p.score,
			p.apy_change_1h, p.apy_change_24h, p.apy_change_7d
		FROM pools p
		WHERE p.apy_change_24h > $1
	`
	args := []interface{}{minGrowth}
	argCount := 1

	if chain != "" {
		argCount++
		query += fmt.Sprintf(" AND p.chain = $%d", argCount)
		args = append(args, chain)
	}

	query += " ORDER BY p.apy_change_24h DESC"

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query trending pools: %w", err)
	}
	defer rows.Close()

	trending := make([]models.TrendingPool, 0)
	for rows.Next() {
		var pool models.Pool
		var change1h, change24h, change7d decimal.Decimal

		err := rows.Scan(
			&pool.ID, &pool.Chain, &pool.Protocol, &pool.Symbol,
			&pool.TVL, &pool.APY, &pool.APYBase, &pool.APYReward, &pool.Score,
			&change1h, &change24h, &change7d,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trending pool: %w", err)
		}

		trending = append(trending, models.TrendingPool{
			Pool:         &pool,
			APYGrowth1H:  change1h,
			APYGrowth24H: change24h,
			APYGrowth7D:  change7d,
			TrendScore:   change24h, // Simple trend score based on 24h growth
		})
	}

	return trending, nil
}

// =============================================================================
// Stats Operations
// =============================================================================

// ListChains returns all chains with aggregated statistics
func (r *Repository) ListChains(ctx context.Context) ([]models.Chain, error) {
	query := `
		SELECT
			chain,
			COUNT(*) as pool_count,
			SUM(tvl) as total_tvl,
			AVG(apy) as average_apy,
			MAX(apy) as max_apy
		FROM pools
		GROUP BY chain
		ORDER BY total_tvl DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query chains: %w", err)
	}
	defer rows.Close()

	chains := make([]models.Chain, 0)
	for rows.Next() {
		var c models.Chain
		err := rows.Scan(
			&c.Name, &c.PoolCount, &c.TotalTVL, &c.AverageAPY, &c.MaxAPY,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chain: %w", err)
		}
		c.DisplayName = c.Name // Can be mapped to human-readable names
		chains = append(chains, c)
	}

	return chains, nil
}

// ListProtocols returns protocols with aggregated statistics
func (r *Repository) ListProtocols(ctx context.Context, filter models.ProtocolFilter) ([]models.Protocol, int64, error) {
	query := `
		SELECT
			protocol,
			array_agg(DISTINCT chain) as chains,
			COUNT(*) as pool_count,
			SUM(tvl) as total_tvl,
			AVG(apy) as average_apy,
			MAX(apy) as max_apy
		FROM pools
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(DISTINCT protocol) FROM pools WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	if filter.Chain != "" {
		argCount++
		query += fmt.Sprintf(" AND chain = $%d", argCount)
		countQuery += fmt.Sprintf(" AND chain = $%d", argCount)
		args = append(args, filter.Chain)
	}

	query += " GROUP BY protocol"

	// Get count
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count protocols: %w", err)
	}

	// Add sorting
	sortColumn := "total_tvl"
	switch filter.SortBy {
	case "poolCount":
		sortColumn = "pool_count"
	case "apy":
		sortColumn = "average_apy"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filter.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filter.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query protocols: %w", err)
	}
	defer rows.Close()

	protocols := make([]models.Protocol, 0)
	for rows.Next() {
		var p models.Protocol
		err := rows.Scan(
			&p.Name, &p.Chains, &p.PoolCount, &p.TotalTVL, &p.AverageAPY, &p.MaxAPY,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan protocol: %w", err)
		}
		p.DisplayName = p.Name
		protocols = append(protocols, p)
	}

	return protocols, total, nil
}

// GetPlatformStats returns overall platform statistics
func (r *Repository) GetPlatformStats(ctx context.Context) (*models.PlatformStats, error) {
	stats := &models.PlatformStats{
		TVLByChain:   make(map[string]decimal.Decimal),
		PoolsByChain: make(map[string]int),
	}

	// Get overall stats
	query := `
		SELECT
			COUNT(*) as total_pools,
			COALESCE(SUM(tvl), 0) as total_tvl,
			COALESCE(AVG(apy), 0) as average_apy,
			COALESCE(MAX(apy), 0) as max_apy,
			COUNT(DISTINCT chain) as total_chains,
			COUNT(DISTINCT protocol) as total_protocols
		FROM pools
	`
	err := r.pool.QueryRow(ctx, query).Scan(
		&stats.TotalPools, &stats.TotalTVL, &stats.AverageAPY,
		&stats.MaxAPY, &stats.TotalChains, &stats.TotalProtocols,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform stats: %w", err)
	}

	// Get active opportunities count
	var activeOpps int
	err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM opportunities WHERE is_active = true").Scan(&activeOpps)
	if err == nil {
		stats.ActiveOpportunities = activeOpps
	}

	// Get TVL by chain
	chainQuery := `
		SELECT chain, SUM(tvl) as tvl, COUNT(*) as pool_count
		FROM pools
		GROUP BY chain
	`
	rows, err := r.pool.Query(ctx, chainQuery)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var chain string
			var tvl decimal.Decimal
			var count int
			if err := rows.Scan(&chain, &tvl, &count); err == nil {
				stats.TVLByChain[chain] = tvl
				stats.PoolsByChain[chain] = count
			}
		}
	}

	// Get APY distribution
	distQuery := `
		SELECT
			COUNT(*) FILTER (WHERE apy >= 0 AND apy < 1) as range_0_1,
			COUNT(*) FILTER (WHERE apy >= 1 AND apy < 5) as range_1_5,
			COUNT(*) FILTER (WHERE apy >= 5 AND apy < 10) as range_5_10,
			COUNT(*) FILTER (WHERE apy >= 10 AND apy < 25) as range_10_25,
			COUNT(*) FILTER (WHERE apy >= 25 AND apy < 50) as range_25_50,
			COUNT(*) FILTER (WHERE apy >= 50 AND apy < 100) as range_50_100,
			COUNT(*) FILTER (WHERE apy >= 100) as range_100_plus
		FROM pools
	`
	err = r.pool.QueryRow(ctx, distQuery).Scan(
		&stats.APYDistribution.Range0to1,
		&stats.APYDistribution.Range1to5,
		&stats.APYDistribution.Range5to10,
		&stats.APYDistribution.Range10to25,
		&stats.APYDistribution.Range25to50,
		&stats.APYDistribution.Range50to100,
		&stats.APYDistribution.Range100Plus,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get APY distribution")
	}

	stats.LastUpdated = time.Now().UTC().Format(time.RFC3339)

	return stats, nil
}

// =============================================================================
// Opportunity Write Operations
// =============================================================================

// UpsertOpportunity inserts or updates an opportunity
func (r *Repository) UpsertOpportunity(ctx context.Context, opp *models.Opportunity) error {
	query := `
		INSERT INTO opportunities (
			id, type, title, description, source_pool_id, target_pool_id,
			pool_id, asset, chain, apy_difference, apy_growth, current_apy,
			potential_profit, tvl, risk_level, score, is_active,
			detected_at, last_seen_at, expires_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			current_apy = EXCLUDED.current_apy,
			potential_profit = EXCLUDED.potential_profit,
			tvl = EXCLUDED.tvl,
			score = EXCLUDED.score,
			is_active = EXCLUDED.is_active,
			last_seen_at = EXCLUDED.last_seen_at,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		opp.ID, opp.Type, opp.Title, opp.Description,
		opp.SourcePoolID, opp.TargetPoolID, opp.PoolID,
		opp.Asset, opp.Chain, opp.APYDifference, opp.APYGrowth,
		opp.CurrentAPY, opp.PotentialProfit, opp.TVL, opp.RiskLevel,
		opp.Score, opp.IsActive, opp.DetectedAt, opp.LastSeenAt,
		opp.ExpiresAt, opp.CreatedAt, opp.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to upsert opportunity: %w", err)
	}

	return nil
}

// DeactivateExpiredOpportunities marks expired opportunities as inactive
func (r *Repository) DeactivateExpiredOpportunities(ctx context.Context) error {
	query := `
		UPDATE opportunities
		SET is_active = false, updated_at = NOW()
		WHERE is_active = true AND expires_at < NOW()
	`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to deactivate expired opportunities: %w", err)
	}

	return nil
}

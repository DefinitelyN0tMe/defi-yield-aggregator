// Package redis provides Redis caching operations for the DeFi Yield Aggregator.
// It handles caching of frequently accessed data and pub/sub for real-time updates.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Cache key prefixes
const (
	PrefixPool          = "pool:"
	PrefixPools         = "pools:"
	PrefixOpportunities = "opportunities:"
	PrefixTrending      = "trending:"
	PrefixChains        = "chains"
	PrefixProtocols     = "protocols:"
	PrefixStats         = "stats"
	PrefixPrices        = "prices:"
)

// Pub/Sub channels
const (
	ChannelPoolUpdates       = "pool_updates"
	ChannelOpportunityAlerts = "opportunity_alerts"
)

// Repository handles all Redis operations
type Repository struct {
	client *redis.Client
}

// NewRepository creates a new Redis repository
func NewRepository(ctx context.Context, cfg config.RedisConfig) (*Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Repository{client: client}, nil
}

// Close closes the Redis connection
func (r *Repository) Close() error {
	return r.client.Close()
}

// Ping checks if Redis connection is alive
func (r *Repository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Client returns the underlying Redis client for advanced operations
func (r *Repository) Client() *redis.Client {
	return r.client
}

// =============================================================================
// Pool Cache Operations
// =============================================================================

// GetPool retrieves a cached pool by ID
func (r *Repository) GetPool(ctx context.Context, id string) (*models.Pool, error) {
	key := PrefixPool + id
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get pool from cache: %w", err)
	}

	var pool models.Pool
	if err := json.Unmarshal(data, &pool); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pool: %w", err)
	}

	return &pool, nil
}

// SetPool caches a pool with TTL in seconds
func (r *Repository) SetPool(ctx context.Context, pool *models.Pool, ttlSeconds int) error {
	key := PrefixPool + pool.ID
	data, err := json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool: %w", err)
	}

	return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// GetPoolsCache retrieves cached pool list response
func (r *Repository) GetPoolsCache(ctx context.Context, cacheKey string) (*models.PoolListResponse, error) {
	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response models.PoolListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetPoolsCache caches a pool list response
func (r *Repository) SetPoolsCache(ctx context.Context, cacheKey string, response *models.PoolListResponse, ttlSeconds int) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// SetMultiplePools caches multiple pools at once using pipeline
func (r *Repository) SetMultiplePools(ctx context.Context, pools []models.Pool, ttlSeconds int) error {
	pipe := r.client.Pipeline()

	for _, pool := range pools {
		key := PrefixPool + pool.ID
		data, err := json.Marshal(pool)
		if err != nil {
			log.Warn().Str("pool_id", pool.ID).Err(err).Msg("Failed to marshal pool")
			continue
		}
		pipe.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// =============================================================================
// Opportunity Cache Operations
// =============================================================================

// GetOpportunitiesCache retrieves cached opportunities
func (r *Repository) GetOpportunitiesCache(ctx context.Context, cacheKey string) (*models.OpportunityListResponse, error) {
	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response models.OpportunityListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetOpportunitiesCache caches opportunities
func (r *Repository) SetOpportunitiesCache(ctx context.Context, cacheKey string, response *models.OpportunityListResponse, ttlSeconds int) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// GetTrendingCache retrieves cached trending pools
func (r *Repository) GetTrendingCache(ctx context.Context, cacheKey string) ([]models.TrendingPool, error) {
	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var trending []models.TrendingPool
	if err := json.Unmarshal(data, &trending); err != nil {
		return nil, err
	}

	return trending, nil
}

// SetTrendingCache caches trending pools
func (r *Repository) SetTrendingCache(ctx context.Context, cacheKey string, trending []models.TrendingPool, ttlSeconds int) error {
	data, err := json.Marshal(trending)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// =============================================================================
// Stats Cache Operations
// =============================================================================

// GetChainsCache retrieves cached chains
func (r *Repository) GetChainsCache(ctx context.Context) (*models.ChainListResponse, error) {
	data, err := r.client.Get(ctx, PrefixChains).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response models.ChainListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetChainsCache caches chains
func (r *Repository) SetChainsCache(ctx context.Context, response *models.ChainListResponse, ttlSeconds int) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, PrefixChains, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// GetProtocolsCache retrieves cached protocols
func (r *Repository) GetProtocolsCache(ctx context.Context, cacheKey string) (*models.ProtocolListResponse, error) {
	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var response models.ProtocolListResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetProtocolsCache caches protocols
func (r *Repository) SetProtocolsCache(ctx context.Context, cacheKey string, response *models.ProtocolListResponse, ttlSeconds int) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, cacheKey, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// GetStatsCache retrieves cached platform stats
func (r *Repository) GetStatsCache(ctx context.Context) (*models.PlatformStats, error) {
	data, err := r.client.Get(ctx, PrefixStats).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var stats models.PlatformStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// SetStatsCache caches platform stats
func (r *Repository) SetStatsCache(ctx context.Context, stats *models.PlatformStats, ttlSeconds int) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, PrefixStats, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// =============================================================================
// Price Cache Operations (for CoinGecko data)
// =============================================================================

// GetTokenPrice retrieves a cached token price
func (r *Repository) GetTokenPrice(ctx context.Context, tokenID string) (float64, error) {
	key := PrefixPrices + tokenID
	price, err := r.client.Get(ctx, key).Float64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return price, nil
}

// SetTokenPrice caches a token price
func (r *Repository) SetTokenPrice(ctx context.Context, tokenID string, price float64, ttlSeconds int) error {
	key := PrefixPrices + tokenID
	return r.client.Set(ctx, key, price, time.Duration(ttlSeconds)*time.Second).Err()
}

// SetMultipleTokenPrices caches multiple token prices using pipeline
func (r *Repository) SetMultipleTokenPrices(ctx context.Context, prices map[string]float64, ttlSeconds int) error {
	pipe := r.client.Pipeline()

	for tokenID, price := range prices {
		key := PrefixPrices + tokenID
		pipe.Set(ctx, key, price, time.Duration(ttlSeconds)*time.Second)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// =============================================================================
// Pub/Sub Operations for Real-Time Updates
// =============================================================================

// PublishPoolUpdate publishes a pool update to subscribers
func (r *Repository) PublishPoolUpdate(ctx context.Context, pool *models.Pool) error {
	data, err := json.Marshal(pool)
	if err != nil {
		return fmt.Errorf("failed to marshal pool for publish: %w", err)
	}

	return r.client.Publish(ctx, ChannelPoolUpdates, data).Err()
}

// PublishOpportunityAlert publishes a new opportunity alert
func (r *Repository) PublishOpportunityAlert(ctx context.Context, opportunity *models.Opportunity) error {
	data, err := json.Marshal(opportunity)
	if err != nil {
		return fmt.Errorf("failed to marshal opportunity for publish: %w", err)
	}

	return r.client.Publish(ctx, ChannelOpportunityAlerts, data).Err()
}

// SubscribePoolUpdates returns a channel for pool update events
func (r *Repository) SubscribePoolUpdates(ctx context.Context) *redis.PubSub {
	return r.client.Subscribe(ctx, ChannelPoolUpdates)
}

// SubscribeOpportunityAlerts returns a channel for opportunity alert events
func (r *Repository) SubscribeOpportunityAlerts(ctx context.Context) *redis.PubSub {
	return r.client.Subscribe(ctx, ChannelOpportunityAlerts)
}

// =============================================================================
// Cache Invalidation
// =============================================================================

// InvalidatePoolCache removes a pool from cache
func (r *Repository) InvalidatePoolCache(ctx context.Context, id string) error {
	return r.client.Del(ctx, PrefixPool+id).Err()
}

// InvalidateAllPoolsCache removes all cached pool lists
func (r *Repository) InvalidateAllPoolsCache(ctx context.Context) error {
	iter := r.client.Scan(ctx, 0, PrefixPools+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			log.Warn().Str("key", iter.Val()).Err(err).Msg("Failed to delete cache key")
		}
	}
	return iter.Err()
}

// InvalidateStatsCache removes all cached stats
func (r *Repository) InvalidateStatsCache(ctx context.Context) error {
	keys := []string{PrefixStats, PrefixChains}
	return r.client.Del(ctx, keys...).Err()
}

// Package defillama provides a client for the DeFiLlama Yields API.
// DeFiLlama aggregates yield data from 2000+ pools across multiple chains.
// API Docs: https://defillama.com/docs/api
package defillama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"golang.org/x/time/rate"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Pool represents a yield pool from DeFiLlama API response
type Pool struct {
	Chain            string   `json:"chain"`
	Project          string   `json:"project"`
	Symbol           string   `json:"symbol"`
	TVLUsd           float64  `json:"tvlUsd"`
	APYBase          float64  `json:"apyBase"`
	APYReward        float64  `json:"apyReward"`
	APY              float64  `json:"apy"`
	RewardTokens     []string `json:"rewardTokens"`
	Pool             string   `json:"pool"`               // Pool ID
	APYPct1D         float64  `json:"apyPct1D"`           // APY change 1 day
	APYPct7D         float64  `json:"apyPct7D"`           // APY change 7 days
	APYPct30D        float64  `json:"apyPct30D"`          // APY change 30 days
	Stablecoin       bool     `json:"stablecoin"`
	ILRisk           string   `json:"ilRisk"`             // Impermanent loss risk
	Exposure         string   `json:"exposure"`           // single, multi
	PredictedClass   string   `json:"predictedClass"`
	APYBase7D        float64  `json:"apyBase7d"`
	APYMean30D       float64  `json:"apyMean30d"`
	VolumeUSD1D      float64  `json:"volumeUsd1d"`
	VolumeUSD7D      float64  `json:"volumeUsd7d"`
	IL7D             float64  `json:"il7d"`
	UnderlyingTokens []string `json:"underlyingTokens"`
	PoolMeta         string   `json:"poolMeta"`
}

// PoolsResponse represents the API response from /pools endpoint
type PoolsResponse struct {
	Status string `json:"status"`
	Data   []Pool `json:"data"`
}

// Client is the DeFiLlama API client
type Client struct {
	baseURL     string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// NewClient creates a new DeFiLlama API client with rate limiting
func NewClient(cfg config.DeFiLlamaConfig) *Client {
	// Calculate rate limiter: requests per minute -> requests per second
	rps := float64(cfg.RateLimit) / 60.0

	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Allow burst of 10 requests, then rate limit
		rateLimiter: rate.NewLimiter(rate.Limit(rps), 10),
	}
}

// FetchPools retrieves all yield pools from DeFiLlama
func (c *Client) FetchPools(ctx context.Context) ([]Pool, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	url := c.baseURL + "/pools"
	log.Debug().Str("url", url).Msg("Fetching pools from DeFiLlama")

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "DeFiYieldAggregator/1.0")

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Msg("DeFiLlama request failed, retrying...")

			if attempt == maxRetries {
				return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
			}

			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
				continue
			}
		}
		break
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var poolsResp PoolsResponse
	if err := json.NewDecoder(resp.Body).Decode(&poolsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info().
		Int("count", len(poolsResp.Data)).
		Msg("Successfully fetched pools from DeFiLlama")

	return poolsResp.Data, nil
}

// FetchPool retrieves a specific pool by ID
func (c *Client) FetchPool(ctx context.Context, poolID string) (*Pool, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	url := fmt.Sprintf("%s/chart/%s", c.baseURL, poolID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var pool Pool
	if err := json.NewDecoder(resp.Body).Decode(&pool); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pool, nil
}

// ToPoolModel converts a DeFiLlama Pool to our internal Pool model
func ToPoolModel(p Pool) models.Pool {
	now := time.Now().UTC()

	return models.Pool{
		ID:               p.Pool,
		Chain:            p.Chain,
		Protocol:         p.Project,
		Symbol:           p.Symbol,
		TVL:              decimal.NewFromFloat(p.TVLUsd),
		APY:              decimal.NewFromFloat(p.APY),
		APYBase:          decimal.NewFromFloat(p.APYBase),
		APYReward:        decimal.NewFromFloat(p.APYReward),
		RewardTokens:     p.RewardTokens,
		UnderlyingTokens: p.UnderlyingTokens,
		PoolMeta:         p.PoolMeta,
		IL7D:             decimal.NewFromFloat(p.IL7D),
		APYMean30D:       decimal.NewFromFloat(p.APYMean30D),
		VolumeUSD1D:      decimal.NewFromFloat(p.VolumeUSD1D),
		VolumeUSD7D:      decimal.NewFromFloat(p.VolumeUSD7D),
		APYChange1H:      decimal.Zero, // Not provided by API, calculated later
		APYChange24H:     decimal.NewFromFloat(p.APYPct1D),
		APYChange7D:      decimal.NewFromFloat(p.APYPct7D),
		StableCoin:       p.Stablecoin,
		Exposure:         p.Exposure,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

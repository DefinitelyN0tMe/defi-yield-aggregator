// Package models contains data structures used throughout the application.
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Pool represents a DeFi yield farming pool
type Pool struct {
	ID              string          `json:"id" db:"id"`                             // Unique identifier (from DeFiLlama)
	Chain           string          `json:"chain" db:"chain"`                       // Blockchain network (ethereum, bsc, polygon, etc.)
	Protocol        string          `json:"protocol" db:"protocol"`                 // Protocol name (aave-v3, compound, curve, etc.)
	Symbol          string          `json:"symbol" db:"symbol"`                     // Pool symbol/name (USDC, ETH-USDC, etc.)
	TVL             decimal.Decimal `json:"tvl" db:"tvl"`                           // Total Value Locked in USD
	APY             decimal.Decimal `json:"apy" db:"apy"`                           // Current Annual Percentage Yield
	APYBase         decimal.Decimal `json:"apyBase" db:"apy_base"`                  // Base APY (from lending/trading fees)
	APYReward       decimal.Decimal `json:"apyReward" db:"apy_reward"`              // Reward APY (from token incentives)
	RewardTokens    []string        `json:"rewardTokens" db:"reward_tokens"`        // Tokens given as rewards
	UnderlyingTokens []string       `json:"underlyingTokens" db:"underlying_tokens"` // Underlying assets in the pool
	PoolMeta        string          `json:"poolMeta" db:"pool_meta"`                // Additional metadata
	IL7D            decimal.Decimal `json:"il7d" db:"il_7d"`                        // 7-day impermanent loss
	APYMean30D      decimal.Decimal `json:"apyMean30d" db:"apy_mean_30d"`           // 30-day average APY
	VolumeUSD1D     decimal.Decimal `json:"volumeUsd1d" db:"volume_usd_1d"`         // 24h trading volume in USD
	VolumeUSD7D     decimal.Decimal `json:"volumeUsd7d" db:"volume_usd_7d"`         // 7-day trading volume in USD

	// Calculated fields
	Score           decimal.Decimal `json:"score" db:"score"`                       // Risk-adjusted opportunity score
	APYChange1H     decimal.Decimal `json:"apyChange1h" db:"apy_change_1h"`         // APY change in last hour
	APYChange24H    decimal.Decimal `json:"apyChange24h" db:"apy_change_24h"`       // APY change in last 24 hours
	APYChange7D     decimal.Decimal `json:"apyChange7d" db:"apy_change_7d"`         // APY change in last 7 days

	// Metadata
	StableCoin      bool            `json:"stablecoin" db:"stablecoin"`             // Is this a stablecoin pool?
	Exposure        string          `json:"exposure" db:"exposure"`                 // Exposure type (single, multi, etc.)

	// Timestamps
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
}

// PoolFilter defines filtering options for pool queries
type PoolFilter struct {
	Chain       string          `query:"chain"`       // Filter by blockchain
	Protocol    string          `query:"protocol"`    // Filter by protocol
	Symbol      string          `query:"symbol"`      // Filter by symbol (partial match)
	MinAPY      decimal.Decimal `query:"minApy"`      // Minimum APY threshold
	MaxAPY      decimal.Decimal `query:"maxApy"`      // Maximum APY threshold
	MinTVL      decimal.Decimal `query:"minTvl"`      // Minimum TVL threshold
	MaxTVL      decimal.Decimal `query:"maxTvl"`      // Maximum TVL threshold
	MinScore    decimal.Decimal `query:"minScore"`    // Minimum score threshold
	StableCoin  *bool           `query:"stablecoin"`  // Filter stablecoin pools
	SortBy      string          `query:"sortBy"`      // Sort field (apy, tvl, score)
	SortOrder   string          `query:"sortOrder"`   // Sort direction (asc, desc)
	Limit       int             `query:"limit"`       // Pagination limit
	Offset      int             `query:"offset"`      // Pagination offset
}

// PoolListResponse is the API response for listing pools
type PoolListResponse struct {
	Data       []Pool `json:"data"`
	Total      int64  `json:"total"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	HasMore    bool   `json:"hasMore"`
}

// HistoricalAPY represents a historical APY data point
type HistoricalAPY struct {
	PoolID    string          `json:"poolId" db:"pool_id"`
	Timestamp time.Time       `json:"timestamp" db:"timestamp"`
	APY       decimal.Decimal `json:"apy" db:"apy"`
	TVL       decimal.Decimal `json:"tvl" db:"tvl"`
	APYBase   decimal.Decimal `json:"apyBase" db:"apy_base"`
	APYReward decimal.Decimal `json:"apyReward" db:"apy_reward"`
}

// PoolHistoryRequest defines the time range for historical data
type PoolHistoryRequest struct {
	Period string `query:"period"` // 1h, 24h, 7d, 30d
}

// PoolHistoryResponse is the API response for pool history
type PoolHistoryResponse struct {
	PoolID    string          `json:"poolId"`
	Period    string          `json:"period"`
	DataPoints []HistoricalAPY `json:"dataPoints"`
}

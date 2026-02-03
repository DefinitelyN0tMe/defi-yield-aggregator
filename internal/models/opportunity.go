package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// OpportunityType defines the type of yield opportunity
type OpportunityType string

const (
	// OpportunityTypeYieldGap represents arbitrage between same asset on different protocols
	OpportunityTypeYieldGap OpportunityType = "yield-gap"
	// OpportunityTypeTrending represents pools with rapidly increasing APY
	OpportunityTypeTrending OpportunityType = "trending"
	// OpportunityTypeHighScore represents pools with high risk-adjusted scores
	OpportunityTypeHighScore OpportunityType = "high-score"
)

// RiskLevel categorizes opportunity risk
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

// Opportunity represents a detected yield farming opportunity
type Opportunity struct {
	ID               string           `json:"id" db:"id"`
	Type             OpportunityType  `json:"type" db:"type"`
	Title            string           `json:"title" db:"title"`
	Description      string           `json:"description" db:"description"`

	// For yield-gap opportunities
	SourcePoolID     string           `json:"sourcePoolId,omitempty" db:"source_pool_id"`
	TargetPoolID     string           `json:"targetPoolId,omitempty" db:"target_pool_id"`
	SourcePool       *Pool            `json:"sourcePool,omitempty" db:"-"`
	TargetPool       *Pool            `json:"targetPool,omitempty" db:"-"`

	// For trending/high-score opportunities
	PoolID           string           `json:"poolId,omitempty" db:"pool_id"`
	Pool             *Pool            `json:"pool,omitempty" db:"-"`

	// Metrics
	Asset            string           `json:"asset" db:"asset"`                     // Base asset (USDC, ETH, etc.)
	Chain            string           `json:"chain" db:"chain"`
	APYDifference    decimal.Decimal  `json:"apyDifference" db:"apy_difference"`    // For yield-gap
	APYGrowth        decimal.Decimal  `json:"apyGrowth" db:"apy_growth"`            // For trending (percentage)
	CurrentAPY       decimal.Decimal  `json:"currentApy" db:"current_apy"`
	PotentialProfit  decimal.Decimal  `json:"potentialProfit" db:"potential_profit"` // Estimated profit in %
	TVL              decimal.Decimal  `json:"tvl" db:"tvl"`                         // Combined or single pool TVL

	// Risk assessment
	RiskLevel        RiskLevel        `json:"riskLevel" db:"risk_level"`
	Score            decimal.Decimal  `json:"score" db:"score"`

	// Status
	IsActive         bool             `json:"isActive" db:"is_active"`
	DetectedAt       time.Time        `json:"detectedAt" db:"detected_at"`
	LastSeenAt       time.Time        `json:"lastSeenAt" db:"last_seen_at"`
	ExpiresAt        time.Time        `json:"expiresAt" db:"expires_at"`

	// Metadata
	CreatedAt        time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time        `json:"updatedAt" db:"updated_at"`
}

// OpportunityFilter defines filtering options for opportunity queries
type OpportunityFilter struct {
	Type        OpportunityType `query:"type"`
	RiskLevel   RiskLevel       `query:"riskLevel"`
	Chain       string          `query:"chain"`
	Asset       string          `query:"asset"`
	MinProfit   decimal.Decimal `query:"minProfit"`
	MinScore    decimal.Decimal `query:"minScore"`
	ActiveOnly  bool            `query:"activeOnly"`
	SortBy      string          `query:"sortBy"`      // profit, score, apy, detectedAt
	SortOrder   string          `query:"sortOrder"`   // asc, desc
	Limit       int             `query:"limit"`
	Offset      int             `query:"offset"`
}

// OpportunityListResponse is the API response for listing opportunities
type OpportunityListResponse struct {
	Data    []Opportunity `json:"data"`
	Total   int64         `json:"total"`
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	HasMore bool          `json:"hasMore"`
}

// TrendingPool represents a pool with significant APY growth
type TrendingPool struct {
	Pool         *Pool           `json:"pool"`
	APYGrowth1H  decimal.Decimal `json:"apyGrowth1h"`  // % growth in 1 hour
	APYGrowth24H decimal.Decimal `json:"apyGrowth24h"` // % growth in 24 hours
	APYGrowth7D  decimal.Decimal `json:"apyGrowth7d"`  // % growth in 7 days
	TrendScore   decimal.Decimal `json:"trendScore"`   // Composite trend score
}

// YieldGap represents an arbitrage opportunity between two pools
type YieldGap struct {
	Asset           string          `json:"asset"`
	LowYieldPool    *Pool           `json:"lowYieldPool"`
	HighYieldPool   *Pool           `json:"highYieldPool"`
	APYDifference   decimal.Decimal `json:"apyDifference"`
	PotentialProfit decimal.Decimal `json:"potentialProfit"`
	Chains          []string        `json:"chains"`
}

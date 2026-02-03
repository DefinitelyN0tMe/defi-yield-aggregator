package models

import (
	"github.com/shopspring/decimal"
)

// Chain represents a blockchain network with aggregated statistics
type Chain struct {
	Name         string          `json:"name"`
	DisplayName  string          `json:"displayName"`
	PoolCount    int             `json:"poolCount"`
	TotalTVL     decimal.Decimal `json:"totalTvl"`
	AverageAPY   decimal.Decimal `json:"averageApy"`
	MaxAPY       decimal.Decimal `json:"maxApy"`
	TopProtocols []string        `json:"topProtocols"`
}

// ChainListResponse is the API response for listing chains
type ChainListResponse struct {
	Data  []Chain `json:"data"`
	Total int     `json:"total"`
}

// Protocol represents a DeFi protocol with aggregated statistics
type Protocol struct {
	Name           string          `json:"name"`
	DisplayName    string          `json:"displayName"`
	Category       string          `json:"category"`    // lending, dex, yield, etc.
	Chains         []string        `json:"chains"`      // Supported chains
	PoolCount      int             `json:"poolCount"`
	TotalTVL       decimal.Decimal `json:"totalTvl"`
	AverageAPY     decimal.Decimal `json:"averageApy"`
	MaxAPY         decimal.Decimal `json:"maxApy"`
	Website        string          `json:"website,omitempty"`
	Twitter        string          `json:"twitter,omitempty"`
	SecurityScore  decimal.Decimal `json:"securityScore"` // 0-100
}

// ProtocolFilter defines filtering options for protocol queries
type ProtocolFilter struct {
	Chain    string `query:"chain"`
	Category string `query:"category"`
	SortBy   string `query:"sortBy"`    // tvl, poolCount, apy
	SortOrder string `query:"sortOrder"` // asc, desc
	Limit    int    `query:"limit"`
	Offset   int    `query:"offset"`
}

// ProtocolListResponse is the API response for listing protocols
type ProtocolListResponse struct {
	Data    []Protocol `json:"data"`
	Total   int64      `json:"total"`
	Limit   int        `json:"limit"`
	Offset  int        `json:"offset"`
	HasMore bool       `json:"hasMore"`
}

// PlatformStats represents overall platform statistics
type PlatformStats struct {
	TotalPools          int             `json:"totalPools"`
	TotalTVL            decimal.Decimal `json:"totalTvl"`
	AverageAPY          decimal.Decimal `json:"averageApy"`
	MedianAPY           decimal.Decimal `json:"medianApy"`
	MaxAPY              decimal.Decimal `json:"maxApy"`
	TotalChains         int             `json:"totalChains"`
	TotalProtocols      int             `json:"totalProtocols"`
	ActiveOpportunities int             `json:"activeOpportunities"`
	LastUpdated         string          `json:"lastUpdated"`

	// Distribution data for charts
	TVLByChain          map[string]decimal.Decimal `json:"tvlByChain"`
	PoolsByChain        map[string]int             `json:"poolsByChain"`
	APYDistribution     APYDistribution            `json:"apyDistribution"`
}

// APYDistribution shows how pools are distributed across APY ranges
type APYDistribution struct {
	Range0to1    int `json:"range0to1"`    // 0-1% APY
	Range1to5    int `json:"range1to5"`    // 1-5% APY
	Range5to10   int `json:"range5to10"`   // 5-10% APY
	Range10to25  int `json:"range10to25"`  // 10-25% APY
	Range25to50  int `json:"range25to50"`  // 25-50% APY
	Range50to100 int `json:"range50to100"` // 50-100% APY
	Range100Plus int `json:"range100plus"` // 100%+ APY
}

// HealthCheck represents the health status of the service
type HealthCheck struct {
	Status      string                 `json:"status"` // healthy, degraded, unhealthy
	Version     string                 `json:"version"`
	Uptime      string                 `json:"uptime"`
	Timestamp   string                 `json:"timestamp"`
	Services    map[string]ServiceHealth `json:"services"`
}

// ServiceHealth represents the health of an individual service
type ServiceHealth struct {
	Status    string `json:"status"`    // up, down
	Latency   string `json:"latency"`   // Response time
	Message   string `json:"message,omitempty"`
}

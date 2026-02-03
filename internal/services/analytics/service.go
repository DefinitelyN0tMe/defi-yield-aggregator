// Package analytics provides scoring algorithms for yield opportunity analysis.
// It calculates risk-adjusted scores for pools based on multiple factors.
package analytics

import (
	"math"

	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Chain security ratings (0-100)
// Higher = more secure/established
var chainSecurityRatings = map[string]float64{
	"ethereum":   95,
	"bsc":        75,
	"polygon":    80,
	"arbitrum":   85,
	"optimism":   85,
	"avalanche":  80,
	"fantom":     70,
	"base":       80,
	"gnosis":     75,
	"celo":       70,
	"moonbeam":   65,
	"moonriver":  60,
	"aurora":     65,
	"cronos":     60,
	"harmony":    50, // Had security issues
	"metis":      60,
	"boba":       55,
	"kava":       65,
	"solana":     75,
}

// Service provides analytics and scoring functionality
type Service struct {
	weights config.ScoringConfig
}

// NewService creates a new analytics service
func NewService(weights config.ScoringConfig) *Service {
	return &Service{weights: weights}
}

// CalculateScore computes a risk-adjusted opportunity score for a pool
// The score is a weighted combination of:
// - APY (higher = better)
// - TVL (higher = safer)
// - Stability (lower volatility = safer)
// - Trend (positive trend = better)
//
// Formula:
// score = (apy_weight * normalized_apy) +
//
//	(tvl_weight * normalized_tvl) +
//	(stability_weight * (1 - volatility)) +
//	(trend_weight * normalized_trend)
func (s *Service) CalculateScore(pool *models.Pool) decimal.Decimal {
	// Normalize APY (0-1 scale, capped at 100%)
	// Uses logarithmic scaling for APY since it can vary widely
	apy, _ := pool.APY.Float64()
	normalizedAPY := normalizeAPY(apy)

	// Normalize TVL (0-1 scale, using logarithmic scaling)
	// Higher TVL = safer (more liquidity, harder to manipulate)
	tvl, _ := pool.TVL.Float64()
	normalizedTVL := normalizeTVL(tvl)

	// Calculate stability score (0-1, higher = more stable)
	// Based on APY volatility over 30 days
	apyMean30d, _ := pool.APYMean30D.Float64()
	stability := calculateStability(apy, apyMean30d)

	// Normalize trend (0-1, positive trend = higher score)
	change24h, _ := pool.APYChange24H.Float64()
	normalizedTrend := normalizeTrend(change24h)

	// Apply chain security multiplier
	chainMultiplier := getChainSecurityMultiplier(pool.Chain)

	// Calculate weighted score
	score := (s.weights.APYWeight * normalizedAPY) +
		(s.weights.TVLWeight * normalizedTVL) +
		(s.weights.StabilityWeight * stability) +
		(s.weights.TrendWeight * normalizedTrend)

	// Apply chain security multiplier
	score *= chainMultiplier

	// Scale to 0-100
	score *= 100

	return decimal.NewFromFloat(math.Max(0, math.Min(100, score)))
}

// normalizeAPY converts APY to a 0-1 scale using logarithmic scaling
// This handles the wide range of APYs (0.1% to 1000%+)
func normalizeAPY(apy float64) float64 {
	if apy <= 0 {
		return 0
	}

	// Use log scaling: score increases logarithmically with APY
	// 1% APY -> ~0.25, 10% APY -> ~0.5, 100% APY -> ~0.75, 1000% APY -> ~1.0
	logAPY := math.Log10(apy + 1)

	// Normalize to 0-1 range (assuming max reasonable APY is 10000%)
	maxLogAPY := math.Log10(10001)
	normalized := logAPY / maxLogAPY

	return math.Min(1, normalized)
}

// normalizeTVL converts TVL to a 0-1 scale using logarithmic scaling
func normalizeTVL(tvl float64) float64 {
	if tvl <= 0 {
		return 0
	}

	// Use log scaling
	// $100K TVL -> ~0.4, $1M TVL -> ~0.5, $10M TVL -> ~0.6, $100M TVL -> ~0.7, $1B TVL -> ~0.8
	logTVL := math.Log10(tvl)

	// Normalize: $1K (3) to $10B (10)
	// Map log values 3-10 to 0-1
	minLog := 3.0  // $1,000
	maxLog := 10.0 // $10,000,000,000

	normalized := (logTVL - minLog) / (maxLog - minLog)

	return math.Max(0, math.Min(1, normalized))
}

// calculateStability calculates how stable the APY has been
// Lower deviation from mean = higher stability score
func calculateStability(currentAPY, meanAPY float64) float64 {
	if meanAPY <= 0 {
		return 0.5 // Unknown stability, return neutral
	}

	// Calculate deviation as percentage of mean
	deviation := math.Abs(currentAPY-meanAPY) / meanAPY

	// Convert to stability score (lower deviation = higher score)
	// 0% deviation = 1.0, 50% deviation = 0.5, 100%+ deviation = 0
	stability := 1 - math.Min(1, deviation)

	return stability
}

// normalizeTrend converts APY change percentage to a 0-1 score
func normalizeTrend(change24h float64) float64 {
	// Positive change = higher score, negative = lower
	// Cap at +/- 100% change
	change := math.Max(-100, math.Min(100, change24h))

	// Convert -100 to +100 range to 0 to 1
	normalized := (change + 100) / 200

	return normalized
}

// getChainSecurityMultiplier returns a multiplier based on chain security
func getChainSecurityMultiplier(chain string) float64 {
	rating, ok := chainSecurityRatings[chain]
	if !ok {
		rating = 50 // Unknown chain gets neutral rating
	}

	// Convert 0-100 rating to 0.5-1.0 multiplier
	// This way even low-rated chains don't get zero score
	return 0.5 + (rating / 200)
}

// CalculateYieldGapProfit calculates potential profit from yield gap arbitrage
// This considers:
// - APY difference
// - Gas costs (estimated based on chain)
// - Minimum investment period to be profitable
func (s *Service) CalculateYieldGapProfit(
	lowAPY, highAPY float64,
	tvl float64,
	sourceChain, targetChain string,
) (profit float64, minDays int) {
	apyDiff := highAPY - lowAPY

	if apyDiff <= 0 {
		return 0, 0
	}

	// Estimate gas costs (simplified)
	gasCostUSD := estimateGasCost(sourceChain) + estimateGasCost(targetChain)

	// Calculate minimum investment to cover gas costs in 7 days
	// profit = (investment * apyDiff/100 / 365 * days) - gasCost
	// To break even in 7 days: investment = gasCost * 365 / (apyDiff * 7)

	if apyDiff > 0 {
		minInvestment := gasCostUSD * 365 / (apyDiff * 7)
		minDays = int(math.Ceil(gasCostUSD * 365 / (apyDiff * 10000)))

		// Calculate profit assuming $10,000 investment over 30 days
		investmentAmount := 10000.0
		profit = (investmentAmount * apyDiff / 100 / 365 * 30) - gasCostUSD

		// If can't break even in 30 days with $10K, not a good opportunity
		if profit < 0 || minInvestment > 100000 {
			return 0, 0
		}
	}

	return profit, minDays
}

// estimateGasCost returns estimated gas cost in USD for transactions on a chain
func estimateGasCost(chain string) float64 {
	// Simplified gas cost estimates (in USD)
	// These would ideally be fetched from a gas oracle
	gasCosts := map[string]float64{
		"ethereum":  50.0, // High gas
		"arbitrum":  1.0,
		"optimism":  1.0,
		"polygon":   0.1,
		"bsc":       0.5,
		"avalanche": 0.5,
		"fantom":    0.1,
		"base":      0.5,
		"gnosis":    0.1,
	}

	cost, ok := gasCosts[chain]
	if !ok {
		return 10.0 // Default estimate for unknown chains
	}

	return cost
}

// CalculateRiskLevel determines the risk level of a pool
func (s *Service) CalculateRiskLevel(pool *models.Pool) models.RiskLevel {
	score, _ := pool.Score.Float64()
	tvl, _ := pool.TVL.Float64()
	apy, _ := pool.APY.Float64()

	// High risk indicators:
	// - Very high APY (>100%)
	// - Low TVL (<$100K)
	// - Low score (<30)
	// - Unknown or low-security chain

	riskFactors := 0

	if apy > 100 {
		riskFactors++
	}
	if apy > 500 {
		riskFactors++
	}

	if tvl < 100000 {
		riskFactors++
	}
	if tvl < 10000 {
		riskFactors++
	}

	if score < 30 {
		riskFactors++
	}

	chainRating := chainSecurityRatings[pool.Chain]
	if chainRating < 60 {
		riskFactors++
	}

	switch {
	case riskFactors >= 3:
		return models.RiskLevelHigh
	case riskFactors >= 1:
		return models.RiskLevelMedium
	default:
		return models.RiskLevelLow
	}
}

// DetectAPYAnomaly checks if APY change is significant enough to alert
func (s *Service) DetectAPYAnomaly(pool *models.Pool, threshold float64) bool {
	change24h, _ := pool.APYChange24H.Float64()

	// Check if APY increased by more than threshold percentage
	return change24h > threshold
}

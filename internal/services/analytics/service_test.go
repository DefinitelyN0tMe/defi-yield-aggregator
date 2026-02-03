package analytics

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

func TestCalculateScore(t *testing.T) {
	cfg := config.ScoringConfig{
		APYWeight:       0.35,
		TVLWeight:       0.25,
		StabilityWeight: 0.25,
		TrendWeight:     0.15,
	}

	service := NewService(cfg)

	tests := []struct {
		name     string
		pool     models.Pool
		minScore float64
		maxScore float64
	}{
		{
			name: "high quality pool",
			pool: models.Pool{
				Chain:       "ethereum",
				APY:         decimal.NewFromFloat(5.0),
				TVL:         decimal.NewFromFloat(100000000), // $100M
				APYMean30D:  decimal.NewFromFloat(5.0),
				APYChange24H: decimal.NewFromFloat(0.1),
			},
			minScore: 50,
			maxScore: 100,
		},
		{
			name: "risky pool",
			pool: models.Pool{
				Chain:       "fantom",
				APY:         decimal.NewFromFloat(500.0), // Very high APY
				TVL:         decimal.NewFromFloat(10000), // Low TVL
				APYMean30D:  decimal.NewFromFloat(100.0),
				APYChange24H: decimal.NewFromFloat(50.0),
			},
			minScore: 0,
			maxScore: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.CalculateScore(&tt.pool)
			scoreFloat, _ := score.Float64()

			if scoreFloat < tt.minScore || scoreFloat > tt.maxScore {
				t.Errorf("Score %.2f outside expected range [%.2f, %.2f]",
					scoreFloat, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestCalculateRiskLevel(t *testing.T) {
	cfg := config.ScoringConfig{
		APYWeight:       0.35,
		TVLWeight:       0.25,
		StabilityWeight: 0.25,
		TrendWeight:     0.15,
	}

	service := NewService(cfg)

	tests := []struct {
		name      string
		pool      models.Pool
		wantLevel models.RiskLevel
	}{
		{
			name: "low risk - stable pool",
			pool: models.Pool{
				Chain: "ethereum",
				APY:   decimal.NewFromFloat(5.0),
				TVL:   decimal.NewFromFloat(100000000),
				Score: decimal.NewFromFloat(80),
			},
			wantLevel: models.RiskLevelLow,
		},
		{
			name: "high risk - high apy low tvl",
			pool: models.Pool{
				Chain: "unknown-chain",
				APY:   decimal.NewFromFloat(1000.0),
				TVL:   decimal.NewFromFloat(5000),
				Score: decimal.NewFromFloat(20),
			},
			wantLevel: models.RiskLevelHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.CalculateRiskLevel(&tt.pool)
			if got != tt.wantLevel {
				t.Errorf("Expected risk level %s, got %s", tt.wantLevel, got)
			}
		})
	}
}

func TestNormalizeAPY(t *testing.T) {
	tests := []struct {
		apy      float64
		minNorm  float64
		maxNorm  float64
	}{
		{0, 0, 0.01},
		{1, 0.05, 0.15},   // Adjusted thresholds
		{10, 0.2, 0.35},   // Adjusted thresholds
		{100, 0.45, 0.6},  // Adjusted thresholds
		{1000, 0.65, 0.85},
	}

	for _, tt := range tests {
		norm := normalizeAPY(tt.apy)
		if norm < tt.minNorm || norm > tt.maxNorm {
			t.Errorf("APY %.2f normalized to %.4f, expected [%.2f, %.2f]",
				tt.apy, norm, tt.minNorm, tt.maxNorm)
		}
	}
}

func TestNormalizeTVL(t *testing.T) {
	tests := []struct {
		tvl      float64
		minNorm  float64
		maxNorm  float64
	}{
		{0, 0, 0.01},
		{100000, 0.2, 0.4},      // $100K
		{1000000, 0.3, 0.5},     // $1M
		{100000000, 0.65, 0.8},  // $100M - adjusted
		{1000000000, 0.8, 0.95}, // $1B - adjusted
	}

	for _, tt := range tests {
		norm := normalizeTVL(tt.tvl)
		if norm < tt.minNorm || norm > tt.maxNorm {
			t.Errorf("TVL %.0f normalized to %.4f, expected [%.2f, %.2f]",
				tt.tvl, norm, tt.minNorm, tt.maxNorm)
		}
	}
}

func TestCalculateStability(t *testing.T) {
	tests := []struct {
		currentAPY float64
		meanAPY    float64
		minStab    float64
		maxStab    float64
	}{
		{5.0, 5.0, 0.9, 1.0},   // No deviation
		{5.0, 10.0, 0.4, 0.6},  // 50% deviation
		{10.0, 5.0, 0.0, 0.1},  // 100% deviation
	}

	for _, tt := range tests {
		stab := calculateStability(tt.currentAPY, tt.meanAPY)
		if stab < tt.minStab || stab > tt.maxStab {
			t.Errorf("Stability for APY %.2f/%.2f is %.4f, expected [%.2f, %.2f]",
				tt.currentAPY, tt.meanAPY, stab, tt.minStab, tt.maxStab)
		}
	}
}

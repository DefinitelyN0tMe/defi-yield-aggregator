// Package opportunity provides yield opportunity detection algorithms.
// It identifies yield gaps, trending pools, and high-score opportunities.
package opportunity

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/postgres"
	"github.com/maxjove/defi-yield-aggregator/internal/repository/redis"
	"github.com/maxjove/defi-yield-aggregator/internal/services/analytics"
)

// Service handles opportunity detection and analysis
type Service struct {
	config    config.WorkerConfig
	pgRepo    *postgres.Repository
	redisRepo *redis.Repository
	analytics *analytics.Service
}

// NewService creates a new opportunity detection service
func NewService(
	cfg config.WorkerConfig,
	pg *postgres.Repository,
	redis *redis.Repository,
	analytics *analytics.Service,
) *Service {
	return &Service{
		config:    cfg,
		pgRepo:    pg,
		redisRepo: redis,
		analytics: analytics,
	}
}

// DetectYieldGaps finds yield gap arbitrage opportunities
// This identifies the same asset with different APYs across protocols
func (s *Service) DetectYieldGaps(ctx context.Context) ([]models.Opportunity, error) {
	log.Debug().Msg("Detecting yield gap opportunities")

	// Fetch all pools above minimum TVL
	filter := models.PoolFilter{
		MinTVL: decimal.NewFromFloat(s.config.MinTVLThreshold),
		Limit:  5000,
	}

	pools, _, err := s.pgRepo.ListPools(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pools: %w", err)
	}

	// Group pools by base asset
	assetPools := groupPoolsByAsset(pools)

	opportunities := make([]models.Opportunity, 0)
	now := time.Now().UTC()

	for asset, assetPoolList := range assetPools {
		if len(assetPoolList) < 2 {
			continue // Need at least 2 pools to compare
		}

		// Sort by APY descending
		sort.Slice(assetPoolList, func(i, j int) bool {
			return assetPoolList[i].APY.GreaterThan(assetPoolList[j].APY)
		})

		// Compare highest APY pools with lowest APY pools
		highestPool := assetPoolList[0]
		lowestPool := assetPoolList[len(assetPoolList)-1]

		apyDiff := highestPool.APY.Sub(lowestPool.APY)
		apyDiffFloat, _ := apyDiff.Float64()

		// Check if difference is above threshold
		if apyDiffFloat >= s.config.YieldGapMinProfit {
			highAPY, _ := highestPool.APY.Float64()
			lowAPY, _ := lowestPool.APY.Float64()
			tvl, _ := highestPool.TVL.Float64()

			// Calculate potential profit
			profit, minDays := s.analytics.CalculateYieldGapProfit(
				lowAPY, highAPY, tvl,
				lowestPool.Chain, highestPool.Chain,
			)

			if profit <= 0 {
				continue
			}

			// Determine risk level
			riskLevel := s.analytics.CalculateRiskLevel(&highestPool)

			opp := models.Opportunity{
				ID:              uuid.New().String(),
				Type:            models.OpportunityTypeYieldGap,
				Title:           fmt.Sprintf("%s Yield Gap: %.2f%% difference", asset, apyDiffFloat),
				Description:     fmt.Sprintf("Move %s from %s (%s) at %.2f%% APY to %s (%s) at %.2f%% APY. Potential profit: $%.2f over 30 days (min %d days to break even)", asset, lowestPool.Protocol, lowestPool.Chain, lowAPY, highestPool.Protocol, highestPool.Chain, highAPY, profit, minDays),
				SourcePoolID:    lowestPool.ID,
				TargetPoolID:    highestPool.ID,
				Asset:           asset,
				Chain:           highestPool.Chain, // Target chain
				APYDifference:   apyDiff,
				CurrentAPY:      highestPool.APY,
				PotentialProfit: decimal.NewFromFloat(profit),
				TVL:             highestPool.TVL.Add(lowestPool.TVL),
				RiskLevel:       riskLevel,
				Score:           highestPool.Score,
				IsActive:        true,
				DetectedAt:      now,
				LastSeenAt:      now,
				ExpiresAt:       now.Add(1 * time.Hour), // Opportunities expire after 1 hour
				CreatedAt:       now,
				UpdatedAt:       now,
			}

			opportunities = append(opportunities, opp)
		}
	}

	// Sort by potential profit descending
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].PotentialProfit.GreaterThan(opportunities[j].PotentialProfit)
	})

	// Keep top 50 opportunities
	if len(opportunities) > 50 {
		opportunities = opportunities[:50]
	}

	log.Info().
		Int("count", len(opportunities)).
		Msg("Detected yield gap opportunities")

	return opportunities, nil
}

// DetectTrendingPools finds pools with rapidly increasing APY
func (s *Service) DetectTrendingPools(ctx context.Context) ([]models.Opportunity, error) {
	log.Debug().Msg("Detecting trending pools")

	// Fetch pools with significant APY growth
	trending, err := s.pgRepo.GetTrendingPools(
		ctx,
		"", // All chains
		decimal.NewFromFloat(s.config.APYJumpThreshold),
		100,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending pools: %w", err)
	}

	opportunities := make([]models.Opportunity, 0)
	now := time.Now().UTC()

	for _, tp := range trending {
		if tp.Pool == nil {
			continue
		}

		pool := tp.Pool
		growth24h, _ := tp.APYGrowth24H.Float64()
		apy, _ := pool.APY.Float64()

		// Determine risk level
		riskLevel := s.analytics.CalculateRiskLevel(pool)

		opp := models.Opportunity{
			ID:          uuid.New().String(),
			Type:        models.OpportunityTypeTrending,
			Title:       fmt.Sprintf("Trending: %s on %s (+%.1f%% APY)", pool.Symbol, pool.Protocol, growth24h),
			Description: fmt.Sprintf("%s pool on %s (%s) has seen APY increase from %.2f%% to %.2f%% in the last 24 hours (%.1f%% growth)", pool.Symbol, pool.Protocol, pool.Chain, apy-growth24h, apy, growth24h),
			PoolID:      pool.ID,
			Asset:       pool.Symbol,
			Chain:       pool.Chain,
			APYGrowth:   tp.APYGrowth24H,
			CurrentAPY:  pool.APY,
			TVL:         pool.TVL,
			RiskLevel:   riskLevel,
			Score:       pool.Score,
			IsActive:    true,
			DetectedAt:  now,
			LastSeenAt:  now,
			ExpiresAt:   now.Add(6 * time.Hour), // Trending opportunities last longer
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		opportunities = append(opportunities, opp)
	}

	log.Info().
		Int("count", len(opportunities)).
		Msg("Detected trending pool opportunities")

	return opportunities, nil
}

// DetectHighScorePools finds pools with excellent risk-adjusted scores
func (s *Service) DetectHighScorePools(ctx context.Context) ([]models.Opportunity, error) {
	log.Debug().Msg("Detecting high-score opportunities")

	// Fetch high-scoring pools
	filter := models.PoolFilter{
		MinScore:  decimal.NewFromFloat(70), // Minimum score of 70/100
		MinTVL:    decimal.NewFromFloat(s.config.MinTVLThreshold),
		MinAPY:    decimal.NewFromFloat(s.config.MinAPYThreshold),
		SortBy:    "score",
		SortOrder: "desc",
		Limit:     100,
	}

	pools, _, err := s.pgRepo.ListPools(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch high-score pools: %w", err)
	}

	opportunities := make([]models.Opportunity, 0)
	now := time.Now().UTC()

	for _, pool := range pools {
		score, _ := pool.Score.Float64()
		apy, _ := pool.APY.Float64()
		tvl, _ := pool.TVL.Float64()

		// Determine risk level (should be low for high-score pools)
		riskLevel := s.analytics.CalculateRiskLevel(&pool)

		opp := models.Opportunity{
			ID:          uuid.New().String(),
			Type:        models.OpportunityTypeHighScore,
			Title:       fmt.Sprintf("High Score: %s on %s (%.1f/100)", pool.Symbol, pool.Protocol, score),
			Description: fmt.Sprintf("%s pool on %s (%s) offers %.2f%% APY with $%.0f TVL. Risk-adjusted score: %.1f/100", pool.Symbol, pool.Protocol, pool.Chain, apy, tvl, score),
			PoolID:      pool.ID,
			Asset:       pool.Symbol,
			Chain:       pool.Chain,
			CurrentAPY:  pool.APY,
			TVL:         pool.TVL,
			RiskLevel:   riskLevel,
			Score:       pool.Score,
			IsActive:    true,
			DetectedAt:  now,
			LastSeenAt:  now,
			ExpiresAt:   now.Add(24 * time.Hour), // High-score opportunities are stable
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		opportunities = append(opportunities, opp)
	}

	log.Info().
		Int("count", len(opportunities)).
		Msg("Detected high-score opportunities")

	return opportunities, nil
}

// groupPoolsByAsset groups pools by their primary asset
// This is used for yield gap detection
func groupPoolsByAsset(pools []models.Pool) map[string][]models.Pool {
	groups := make(map[string][]models.Pool)

	for _, pool := range pools {
		// Normalize asset name
		asset := normalizeAsset(pool.Symbol)
		if asset == "" {
			continue
		}

		groups[asset] = append(groups[asset], pool)
	}

	return groups
}

// normalizeAsset extracts and normalizes the primary asset from a pool symbol
func normalizeAsset(symbol string) string {
	// Handle common patterns
	symbol = strings.ToUpper(symbol)

	// Single asset pools
	singleAssets := []string{"USDC", "USDT", "DAI", "FRAX", "LUSD", "BUSD", "ETH", "WETH", "BTC", "WBTC", "MATIC", "AVAX", "FTM", "BNB"}
	for _, asset := range singleAssets {
		if symbol == asset || strings.HasPrefix(symbol, asset+"-") || strings.HasSuffix(symbol, "-"+asset) {
			return asset
		}
	}

	// LP tokens (e.g., "ETH-USDC", "WETH/DAI")
	separators := []string{"-", "/", "_"}
	for _, sep := range separators {
		parts := strings.Split(symbol, sep)
		if len(parts) >= 2 {
			// Return the first token for grouping
			return strings.TrimSpace(parts[0])
		}
	}

	// Return as-is for unknown patterns
	return symbol
}

// RefreshOpportunities updates existing opportunities and deactivates expired ones
func (s *Service) RefreshOpportunities(ctx context.Context) error {
	log.Debug().Msg("Refreshing opportunities")

	// This would update existing opportunities and deactivate expired ones
	// Implementation depends on specific requirements

	return nil
}

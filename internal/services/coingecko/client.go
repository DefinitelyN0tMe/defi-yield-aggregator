// Package coingecko provides a client for the CoinGecko API.
// CoinGecko provides cryptocurrency price data for profitability calculations.
// API Docs: https://www.coingecko.com/api/documentation
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
)

// PriceResponse represents the API response from /simple/price endpoint
// Example: {"ethereum":{"usd":3500.50},"bitcoin":{"usd":45000.00}}
type PriceResponse map[string]map[string]float64

// Client is the CoinGecko API client
type Client struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

// NewClient creates a new CoinGecko API client with rate limiting
func NewClient(cfg config.CoinGeckoConfig) *Client {
	// CoinGecko Demo plan: 30 requests/min
	rps := float64(cfg.RateLimit) / 60.0

	return &Client{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Allow burst of 5 requests, then rate limit
		rateLimiter: rate.NewLimiter(rate.Limit(rps), 5),
	}
}

// FetchPrices retrieves prices for multiple tokens in USD
func (c *Client) FetchPrices(ctx context.Context, tokenIDs []string) (map[string]float64, error) {
	if len(tokenIDs) == 0 {
		return make(map[string]float64), nil
	}

	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Build URL with token IDs
	ids := strings.Join(tokenIDs, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, ids)

	log.Debug().
		Str("url", url).
		Int("token_count", len(tokenIDs)).
		Msg("Fetching prices from CoinGecko")

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "DeFiYieldAggregator/1.0")

	// Add API key if available (for higher rate limits)
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Msg("CoinGecko request failed, retrying...")

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

		// Check for rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Warn().
				Int("attempt", attempt).
				Msg("CoinGecko rate limit hit, waiting...")

			resp.Body.Close()

			if attempt == maxRetries {
				return nil, fmt.Errorf("rate limit exceeded after %d attempts", maxRetries)
			}

			// Wait longer for rate limit
			backoff := time.Duration(attempt*10) * time.Second
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
	var priceResp PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract USD prices
	prices := make(map[string]float64)
	for tokenID, currencies := range priceResp {
		if usdPrice, ok := currencies["usd"]; ok {
			prices[tokenID] = usdPrice
		}
	}

	log.Info().
		Int("count", len(prices)).
		Msg("Successfully fetched prices from CoinGecko")

	return prices, nil
}

// FetchPrice retrieves the price for a single token
func (c *Client) FetchPrice(ctx context.Context, tokenID string) (float64, error) {
	prices, err := c.FetchPrices(ctx, []string{tokenID})
	if err != nil {
		return 0, err
	}

	price, ok := prices[tokenID]
	if !ok {
		return 0, fmt.Errorf("price not found for token: %s", tokenID)
	}

	return price, nil
}

// FetchMarketData retrieves detailed market data for tokens
func (c *Client) FetchMarketData(ctx context.Context, tokenIDs []string) ([]MarketData, error) {
	if len(tokenIDs) == 0 {
		return []MarketData{}, nil
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	ids := strings.Join(tokenIDs, ",")
	url := fmt.Sprintf(
		"%s/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&sparkline=false",
		c.baseURL, ids,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-cg-demo-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var marketData []MarketData
	if err := json.NewDecoder(resp.Body).Decode(&marketData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return marketData, nil
}

// MarketData represents detailed token market data
type MarketData struct {
	ID                           string  `json:"id"`
	Symbol                       string  `json:"symbol"`
	Name                         string  `json:"name"`
	CurrentPrice                 float64 `json:"current_price"`
	MarketCap                    float64 `json:"market_cap"`
	MarketCapRank                int     `json:"market_cap_rank"`
	TotalVolume                  float64 `json:"total_volume"`
	High24H                      float64 `json:"high_24h"`
	Low24H                       float64 `json:"low_24h"`
	PriceChange24H               float64 `json:"price_change_24h"`
	PriceChangePercentage24H     float64 `json:"price_change_percentage_24h"`
	CirculatingSupply            float64 `json:"circulating_supply"`
	TotalSupply                  float64 `json:"total_supply"`
	ATH                          float64 `json:"ath"`
	ATHChangePercentage          float64 `json:"ath_change_percentage"`
	ATL                          float64 `json:"atl"`
	ATLChangePercentage          float64 `json:"atl_change_percentage"`
}

// Common token ID mappings (symbol -> CoinGecko ID)
var TokenIDMap = map[string]string{
	"ETH":   "ethereum",
	"WETH":  "weth",
	"BTC":   "bitcoin",
	"WBTC":  "wrapped-bitcoin",
	"USDC":  "usd-coin",
	"USDT":  "tether",
	"DAI":   "dai",
	"BUSD":  "binance-usd",
	"FRAX":  "frax",
	"LUSD":  "liquity-usd",
	"SUSD":  "nusd",
	"BNB":   "binancecoin",
	"MATIC": "matic-network",
	"AVAX":  "avalanche-2",
	"FTM":   "fantom",
	"ARB":   "arbitrum",
	"OP":    "optimism",
	"CRV":   "curve-dao-token",
	"CVX":   "convex-finance",
	"AAVE":  "aave",
	"COMP":  "compound-governance-token",
	"UNI":   "uniswap",
	"SUSHI": "sushi",
	"MKR":   "maker",
	"SNX":   "havven",
	"YFI":   "yearn-finance",
	"LINK":  "chainlink",
}

// GetTokenID returns the CoinGecko ID for a token symbol
func GetTokenID(symbol string) string {
	symbol = strings.ToUpper(symbol)
	if id, ok := TokenIDMap[symbol]; ok {
		return id
	}
	// Return lowercase symbol as fallback
	return strings.ToLower(symbol)
}

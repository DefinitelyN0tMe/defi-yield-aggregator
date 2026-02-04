// Package elasticsearch provides ElasticSearch operations for the DeFi Yield Aggregator.
// It handles fast searching, filtering, and analytics across pool data.
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/config"
	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Index names
const (
	IndexPools         = "defi_pools"
	IndexOpportunities = "defi_opportunities"
)

// Repository handles all ElasticSearch operations
type Repository struct {
	client *elasticsearch.Client
}

// NewRepository creates a new ElasticSearch repository
func NewRepository(cfg config.ElasticSearchConfig) (*Repository, error) {
	esConfig := elasticsearch.Config{
		Addresses: []string{cfg.URL},
	}

	// Add authentication if configured
	if cfg.Username != "" && cfg.Password != "" {
		esConfig.Username = cfg.Username
		esConfig.Password = cfg.Password
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ElasticSearch client: %w", err)
	}

	return &Repository{client: client}, nil
}

// Ping checks if ElasticSearch connection is alive
func (r *Repository) Ping(ctx context.Context) error {
	res, err := r.client.Ping(r.client.Ping.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to ping ElasticSearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ElasticSearch ping failed: %s", res.Status())
	}

	return nil
}

// CreateIndices creates the necessary indices with mappings
func (r *Repository) CreateIndices(ctx context.Context) error {
	// Create pools index
	if err := r.createPoolsIndex(ctx); err != nil {
		return err
	}

	// Create opportunities index
	if err := r.createOpportunitiesIndex(ctx); err != nil {
		return err
	}

	return nil
}

// createPoolsIndex creates the pools index with proper mappings
func (r *Repository) createPoolsIndex(ctx context.Context) error {
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"lowercase_analyzer": {
						"type": "custom",
						"tokenizer": "standard",
						"filter": ["lowercase"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"chain": {
					"type": "text",
					"analyzer": "lowercase_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"protocol": {
					"type": "text",
					"analyzer": "lowercase_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"symbol": {
					"type": "text",
					"analyzer": "lowercase_analyzer",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				},
				"tvl": { "type": "double" },
				"apy": { "type": "double" },
				"apy_base": { "type": "double" },
				"apy_reward": { "type": "double" },
				"reward_tokens": { "type": "keyword" },
				"underlying_tokens": { "type": "keyword" },
				"pool_meta": { "type": "text" },
				"il_7d": { "type": "double" },
				"apy_mean_30d": { "type": "double" },
				"volume_usd_1d": { "type": "double" },
				"volume_usd_7d": { "type": "double" },
				"score": { "type": "double" },
				"apy_change_1h": { "type": "double" },
				"apy_change_24h": { "type": "double" },
				"apy_change_7d": { "type": "double" },
				"stablecoin": { "type": "boolean" },
				"exposure": { "type": "keyword" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" }
			}
		}
	}`

	res, err := r.client.Indices.Create(
		IndexPools,
		r.client.Indices.Create.WithContext(ctx),
		r.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("failed to create pools index: %w", err)
	}
	defer res.Body.Close()

	// Ignore "already exists" error
	if res.IsError() && !strings.Contains(res.String(), "resource_already_exists_exception") {
		return fmt.Errorf("failed to create pools index: %s", res.String())
	}

	log.Info().Msg("Pools index created/verified")
	return nil
}

// createOpportunitiesIndex creates the opportunities index
func (r *Repository) createOpportunitiesIndex(ctx context.Context) error {
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"type": { "type": "keyword" },
				"title": { "type": "text" },
				"description": { "type": "text" },
				"source_pool_id": { "type": "keyword" },
				"target_pool_id": { "type": "keyword" },
				"pool_id": { "type": "keyword" },
				"asset": { "type": "keyword" },
				"chain": { "type": "keyword" },
				"apy_difference": { "type": "double" },
				"apy_growth": { "type": "double" },
				"current_apy": { "type": "double" },
				"potential_profit": { "type": "double" },
				"tvl": { "type": "double" },
				"risk_level": { "type": "keyword" },
				"score": { "type": "double" },
				"is_active": { "type": "boolean" },
				"detected_at": { "type": "date" },
				"last_seen_at": { "type": "date" },
				"expires_at": { "type": "date" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" }
			}
		}
	}`

	res, err := r.client.Indices.Create(
		IndexOpportunities,
		r.client.Indices.Create.WithContext(ctx),
		r.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("failed to create opportunities index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.String(), "resource_already_exists_exception") {
		return fmt.Errorf("failed to create opportunities index: %s", res.String())
	}

	log.Info().Msg("Opportunities index created/verified")
	return nil
}

// =============================================================================
// Pool Search Operations
// =============================================================================

// SearchPools performs a filtered search on pools
func (r *Repository) SearchPools(ctx context.Context, filter models.PoolFilter) ([]models.Pool, int64, error) {
	// Build ElasticSearch query
	query := buildPoolSearchQuery(filter)

	// Serialize query
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, fmt.Errorf("failed to encode query: %w", err)
	}

	// Execute search
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(IndexPools),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search pools: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	// Parse response
	var result searchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	pools := make([]models.Pool, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var pool models.Pool
		if err := json.Unmarshal(hit.Source, &pool); err != nil {
			log.Warn().Err(err).Str("id", hit.ID).Msg("Failed to unmarshal pool")
			continue
		}
		pools = append(pools, pool)
	}

	return pools, result.Hits.Total.Value, nil
}

// buildPoolSearchQuery builds an ElasticSearch query from filter parameters
func buildPoolSearchQuery(filter models.PoolFilter) map[string]interface{} {
	must := make([]map[string]interface{}, 0)

	// Chain filter (case-insensitive)
	if filter.Chain != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"chain": map[string]interface{}{
					"query":    strings.ToLower(filter.Chain),
					"operator": "and",
				},
			},
		})
	}

	// Protocol filter (case-insensitive)
	if filter.Protocol != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"protocol": map[string]interface{}{
					"query":    strings.ToLower(filter.Protocol),
					"operator": "and",
				},
			},
		})
	}

	// Symbol search (fuzzy match)
	if filter.Symbol != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"symbol": map[string]interface{}{
					"query":     filter.Symbol,
					"fuzziness": "AUTO",
				},
			},
		})
	}

	// General search across multiple fields
	if filter.Search != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     filter.Search,
				"fields":    []string{"symbol^3", "protocol^2", "chain", "pool_meta"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// APY range
	apyRange := make(map[string]interface{})
	if !filter.MinAPY.IsZero() {
		apyRange["gte"], _ = filter.MinAPY.Float64()
	}
	if !filter.MaxAPY.IsZero() {
		apyRange["lte"], _ = filter.MaxAPY.Float64()
	}
	if len(apyRange) > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"apy": apyRange,
			},
		})
	}

	// TVL range
	tvlRange := make(map[string]interface{})
	if !filter.MinTVL.IsZero() {
		tvlRange["gte"], _ = filter.MinTVL.Float64()
	}
	if !filter.MaxTVL.IsZero() {
		tvlRange["lte"], _ = filter.MaxTVL.Float64()
	}
	if len(tvlRange) > 0 {
		must = append(must, map[string]interface{}{
			"range": map[string]interface{}{
				"tvl": tvlRange,
			},
		})
	}

	// Stablecoin filter
	if filter.StableCoin != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"stablecoin": *filter.StableCoin,
			},
		})
	}

	// Build query
	var boolQuery map[string]interface{}
	if len(must) > 0 {
		boolQuery = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		}
	} else {
		boolQuery = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// Build sort
	sortField := "tvl"
	switch filter.SortBy {
	case "apy":
		sortField = "apy"
	case "score":
		sortField = "score"
	}

	sortOrder := "desc"
	if filter.SortOrder == "asc" {
		sortOrder = "asc"
	}

	return map[string]interface{}{
		"query": boolQuery,
		"sort": []map[string]interface{}{
			{
				sortField: map[string]interface{}{
					"order": sortOrder,
				},
			},
		},
		"from": filter.Offset,
		"size": filter.Limit,
	}
}

// =============================================================================
// Index Operations
// =============================================================================

// IndexPool indexes a single pool document
func (r *Repository) IndexPool(ctx context.Context, pool *models.Pool) error {
	doc := poolToDocument(pool)

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal pool: %w", err)
	}

	res, err := r.client.Index(
		IndexPools,
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(pool.ID),
		r.client.Index.WithRefresh("false"), // Don't wait for refresh
	)
	if err != nil {
		return fmt.Errorf("failed to index pool: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}

	return nil
}

// BulkIndexPools indexes multiple pools efficiently
func (r *Repository) BulkIndexPools(ctx context.Context, pools []models.Pool) error {
	if len(pools) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for _, pool := range pools {
		// Action line
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": IndexPools,
				"_id":    pool.ID,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("failed to encode meta: %w", err)
		}

		// Document line
		doc := poolToDocument(&pool)
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return fmt.Errorf("failed to encode document: %w", err)
		}
	}

	res, err := r.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.client.Bulk.WithContext(ctx),
		r.client.Bulk.WithRefresh("false"),
	)
	if err != nil {
		return fmt.Errorf("failed to bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing error: %s", res.String())
	}

	log.Info().Int("count", len(pools)).Msg("Bulk indexed pools")
	return nil
}

// IndexOpportunity indexes a single opportunity
func (r *Repository) IndexOpportunity(ctx context.Context, opp *models.Opportunity) error {
	data, err := json.Marshal(opp)
	if err != nil {
		return fmt.Errorf("failed to marshal opportunity: %w", err)
	}

	res, err := r.client.Index(
		IndexOpportunities,
		bytes.NewReader(data),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(opp.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to index opportunity: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("indexing error: %s", res.String())
	}

	return nil
}

// RefreshIndex forces a refresh of an index
func (r *Repository) RefreshIndex(ctx context.Context, index string) error {
	res, err := r.client.Indices.Refresh(
		r.client.Indices.Refresh.WithContext(ctx),
		r.client.Indices.Refresh.WithIndex(index),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("refresh error: %s", res.String())
	}

	return nil
}

// =============================================================================
// Analytics Operations
// =============================================================================

// GetPoolAggregations returns aggregated pool statistics
func (r *Repository) GetPoolAggregations(ctx context.Context) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"chains": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "chain",
					"size":  100,
				},
				"aggs": map[string]interface{}{
					"total_tvl": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "tvl",
						},
					},
					"avg_apy": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "apy",
						},
					},
				},
			},
			"protocols": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "protocol",
					"size":  100,
				},
				"aggs": map[string]interface{}{
					"total_tvl": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "tvl",
						},
					},
				},
			},
			"total_tvl": map[string]interface{}{
				"sum": map[string]interface{}{
					"field": "tvl",
				},
			},
			"avg_apy": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "apy",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(IndexPools),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// =============================================================================
// Helper Types and Functions
// =============================================================================

// searchResponse represents an ElasticSearch search response
type searchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string          `json:"_id"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// esDocument represents a pool document for ElasticSearch
type esDocument struct {
	ID               string   `json:"id"`
	Chain            string   `json:"chain"`
	Protocol         string   `json:"protocol"`
	Symbol           string   `json:"symbol"`
	TVL              float64  `json:"tvl"`
	APY              float64  `json:"apy"`
	APYBase          float64  `json:"apy_base"`
	APYReward        float64  `json:"apy_reward"`
	RewardTokens     []string `json:"reward_tokens"`
	UnderlyingTokens []string `json:"underlying_tokens"`
	PoolMeta         string   `json:"pool_meta"`
	IL7D             float64  `json:"il_7d"`
	APYMean30D       float64  `json:"apy_mean_30d"`
	VolumeUSD1D      float64  `json:"volume_usd_1d"`
	VolumeUSD7D      float64  `json:"volume_usd_7d"`
	Score            float64  `json:"score"`
	APYChange1H      float64  `json:"apy_change_1h"`
	APYChange24H     float64  `json:"apy_change_24h"`
	APYChange7D      float64  `json:"apy_change_7d"`
	StableCoin       bool     `json:"stablecoin"`
	Exposure         string   `json:"exposure"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

// poolToDocument converts a Pool model to an ElasticSearch document
func poolToDocument(pool *models.Pool) esDocument {
	return esDocument{
		ID:               pool.ID,
		Chain:            pool.Chain,
		Protocol:         pool.Protocol,
		Symbol:           pool.Symbol,
		TVL:              decimalToFloat(pool.TVL),
		APY:              decimalToFloat(pool.APY),
		APYBase:          decimalToFloat(pool.APYBase),
		APYReward:        decimalToFloat(pool.APYReward),
		RewardTokens:     pool.RewardTokens,
		UnderlyingTokens: pool.UnderlyingTokens,
		PoolMeta:         pool.PoolMeta,
		IL7D:             decimalToFloat(pool.IL7D),
		APYMean30D:       decimalToFloat(pool.APYMean30D),
		VolumeUSD1D:      decimalToFloat(pool.VolumeUSD1D),
		VolumeUSD7D:      decimalToFloat(pool.VolumeUSD7D),
		Score:            decimalToFloat(pool.Score),
		APYChange1H:      decimalToFloat(pool.APYChange1H),
		APYChange24H:     decimalToFloat(pool.APYChange24H),
		APYChange7D:      decimalToFloat(pool.APYChange7D),
		StableCoin:       pool.StableCoin,
		Exposure:         pool.Exposure,
		CreatedAt:        pool.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        pool.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func decimalToFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

// Ensure esapi is used to avoid import error
var _ esapi.Response

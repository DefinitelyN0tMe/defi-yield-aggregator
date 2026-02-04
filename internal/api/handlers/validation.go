package handlers

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

// Validation constants
const (
	MaxLimit     = 100
	DefaultLimit = 50
	MaxOffset    = 10000
)

// Valid sort fields for pools
var validPoolSortFields = map[string]bool{
	"apy":        true,
	"tvl":        true,
	"score":      true,
	"updated_at": true,
	"chain":      true,
	"protocol":   true,
}

// Valid sort fields for opportunities
var validOpportunitySortFields = map[string]bool{
	"score":       true,
	"profit":      true,
	"apy":         true,
	"detected_at": true,
}

// Valid time periods
var validPeriods = map[string]bool{
	"1h":  true,
	"24h": true,
	"7d":  true,
	"30d": true,
}

// Valid opportunity types
var validOpportunityTypes = map[string]bool{
	"yield-gap":  true,
	"trending":   true,
	"high-score": true,
}

// Valid risk levels
var validRiskLevels = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

// chainRegex validates chain names
var chainRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// protocolRegex validates protocol names
var protocolRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// ParsePoolFilter parses and validates pool filter parameters
func ParsePoolFilter(c *fiber.Ctx) (models.PoolFilter, []ValidationError) {
	var errors []ValidationError

	filter := models.PoolFilter{
		Chain:     c.Query("chain"),
		Protocol:  c.Query("protocol"),
		Symbol:    c.Query("symbol"),
		Search:    c.Query("search"),
		SortBy:    c.Query("sortBy", "tvl"),
		SortOrder: strings.ToLower(c.Query("sortOrder", "desc")),
		Limit:     c.QueryInt("limit", DefaultLimit),
		Offset:    c.QueryInt("offset", 0),
	}

	// Parse decimal values
	if minApy := c.Query("minApy"); minApy != "" {
		if d, err := decimal.NewFromString(minApy); err != nil {
			errors = append(errors, ValidationError{Field: "minApy", Message: "must be a valid number"})
		} else if d.IsNegative() {
			errors = append(errors, ValidationError{Field: "minApy", Message: "must be non-negative"})
		} else {
			filter.MinAPY = d
		}
	}

	if maxApy := c.Query("maxApy"); maxApy != "" {
		if d, err := decimal.NewFromString(maxApy); err != nil {
			errors = append(errors, ValidationError{Field: "maxApy", Message: "must be a valid number"})
		} else if d.IsNegative() {
			errors = append(errors, ValidationError{Field: "maxApy", Message: "must be non-negative"})
		} else {
			filter.MaxAPY = d
		}
	}

	if minTvl := c.Query("minTvl"); minTvl != "" {
		if d, err := decimal.NewFromString(minTvl); err != nil {
			errors = append(errors, ValidationError{Field: "minTvl", Message: "must be a valid number"})
		} else if d.IsNegative() {
			errors = append(errors, ValidationError{Field: "minTvl", Message: "must be non-negative"})
		} else {
			filter.MinTVL = d
		}
	}

	if maxTvl := c.Query("maxTvl"); maxTvl != "" {
		if d, err := decimal.NewFromString(maxTvl); err != nil {
			errors = append(errors, ValidationError{Field: "maxTvl", Message: "must be a valid number"})
		} else if d.IsNegative() {
			errors = append(errors, ValidationError{Field: "maxTvl", Message: "must be non-negative"})
		} else {
			filter.MaxTVL = d
		}
	}

	if minScore := c.Query("minScore"); minScore != "" {
		if d, err := decimal.NewFromString(minScore); err != nil {
			errors = append(errors, ValidationError{Field: "minScore", Message: "must be a valid number"})
		} else if d.IsNegative() || d.GreaterThan(decimal.NewFromInt(100)) {
			errors = append(errors, ValidationError{Field: "minScore", Message: "must be between 0 and 100"})
		} else {
			filter.MinScore = d
		}
	}

	// Parse stablecoin filter
	if stablecoin := c.Query("stablecoin"); stablecoin != "" {
		val := stablecoin == "true" || stablecoin == "1"
		filter.StableCoin = &val
	}

	// Chain and protocol validation - allow alphanumeric with dashes, underscores, and spaces
	// No strict validation needed as we use case-insensitive matching in the database

	// Validate sort field
	if !validPoolSortFields[filter.SortBy] {
		errors = append(errors, ValidationError{Field: "sortBy", Message: "invalid sort field"})
	}

	// Validate sort order
	if filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		errors = append(errors, ValidationError{Field: "sortOrder", Message: "must be 'asc' or 'desc'"})
	}

	// Validate limit
	if filter.Limit < 1 {
		filter.Limit = DefaultLimit
	} else if filter.Limit > MaxLimit {
		filter.Limit = MaxLimit
	}

	// Validate offset
	if filter.Offset < 0 {
		filter.Offset = 0
	} else if filter.Offset > MaxOffset {
		errors = append(errors, ValidationError{Field: "offset", Message: "exceeds maximum value"})
	}

	// Validate APY range
	if !filter.MinAPY.IsZero() && !filter.MaxAPY.IsZero() && filter.MinAPY.GreaterThan(filter.MaxAPY) {
		errors = append(errors, ValidationError{Field: "minApy", Message: "minApy cannot be greater than maxApy"})
	}

	// Validate TVL range
	if !filter.MinTVL.IsZero() && !filter.MaxTVL.IsZero() && filter.MinTVL.GreaterThan(filter.MaxTVL) {
		errors = append(errors, ValidationError{Field: "minTvl", Message: "minTvl cannot be greater than maxTvl"})
	}

	return filter, errors
}

// ParseOpportunityFilter parses and validates opportunity filter parameters
func ParseOpportunityFilter(c *fiber.Ctx) (models.OpportunityFilter, []ValidationError) {
	var errors []ValidationError

	filter := models.OpportunityFilter{
		Type:       models.OpportunityType(c.Query("type")),
		RiskLevel:  models.RiskLevel(c.Query("riskLevel")),
		Chain:      strings.ToLower(c.Query("chain")),
		Asset:      strings.ToUpper(c.Query("asset")),
		ActiveOnly: c.QueryBool("activeOnly", true),
		SortBy:     c.Query("sortBy", "score"),
		SortOrder:  strings.ToLower(c.Query("sortOrder", "desc")),
		Limit:      c.QueryInt("limit", DefaultLimit),
		Offset:     c.QueryInt("offset", 0),
	}

	// Parse minProfit
	if minProfit := c.Query("minProfit"); minProfit != "" {
		if d, err := decimal.NewFromString(minProfit); err != nil {
			errors = append(errors, ValidationError{Field: "minProfit", Message: "must be a valid number"})
		} else if d.IsNegative() {
			errors = append(errors, ValidationError{Field: "minProfit", Message: "must be non-negative"})
		} else {
			filter.MinProfit = d
		}
	}

	// Parse minScore
	if minScore := c.Query("minScore"); minScore != "" {
		if d, err := decimal.NewFromString(minScore); err != nil {
			errors = append(errors, ValidationError{Field: "minScore", Message: "must be a valid number"})
		} else {
			filter.MinScore = d
		}
	}

	// Validate type
	if filter.Type != "" && !validOpportunityTypes[string(filter.Type)] {
		errors = append(errors, ValidationError{Field: "type", Message: "invalid opportunity type"})
	}

	// Validate risk level
	if filter.RiskLevel != "" && !validRiskLevels[string(filter.RiskLevel)] {
		errors = append(errors, ValidationError{Field: "riskLevel", Message: "invalid risk level"})
	}

	// Validate sort field
	if !validOpportunitySortFields[filter.SortBy] {
		errors = append(errors, ValidationError{Field: "sortBy", Message: "invalid sort field"})
	}

	// Validate sort order
	if filter.SortOrder != "asc" && filter.SortOrder != "desc" {
		errors = append(errors, ValidationError{Field: "sortOrder", Message: "must be 'asc' or 'desc'"})
	}

	// Validate limit
	if filter.Limit < 1 {
		filter.Limit = DefaultLimit
	} else if filter.Limit > MaxLimit {
		filter.Limit = MaxLimit
	}

	// Validate offset
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return filter, errors
}

// ValidatePoolID validates a pool ID
func ValidatePoolID(id string) []ValidationError {
	var errors []ValidationError

	if id == "" {
		errors = append(errors, ValidationError{Field: "id", Message: "pool ID is required"})
	} else if len(id) > 255 {
		errors = append(errors, ValidationError{Field: "id", Message: "pool ID too long"})
	}

	return errors
}

// ValidatePeriod validates a time period parameter
func ValidatePeriod(period string) []ValidationError {
	var errors []ValidationError

	if period != "" && !validPeriods[period] {
		errors = append(errors, ValidationError{Field: "period", Message: "must be one of: 1h, 24h, 7d, 30d"})
	}

	return errors
}

package handlers

import (
	"testing"

	"github.com/maxjove/defi-yield-aggregator/internal/models"
)

func TestParsePoolFilter_Defaults(t *testing.T) {
	// This would require a mock Fiber context
	// For now, test the validation logic directly
}

func TestBuildPoolsCacheKey(t *testing.T) {
	filter := models.PoolFilter{
		Chain:     "ethereum",
		Protocol:  "aave-v3",
		SortBy:    "tvl",
		SortOrder: "desc",
		Limit:     50,
		Offset:    0,
	}

	key := buildPoolsCacheKey(filter)
	expected := "pools:ethereum:aave-v3:tvl:desc:50:0"

	if key != expected {
		t.Errorf("Expected cache key %s, got %s", expected, key)
	}
}

func TestBuildOpportunitiesCacheKey(t *testing.T) {
	filter := models.OpportunityFilter{
		Type:       models.OpportunityTypeYieldGap,
		RiskLevel:  models.RiskLevelLow,
		Chain:      "ethereum",
		SortBy:     "score",
		SortOrder:  "desc",
		ActiveOnly: true,
	}

	key := buildOpportunitiesCacheKey(filter)
	expected := "opportunities:yield-gap:low:ethereum:score:desc:true"

	if key != expected {
		t.Errorf("Expected cache key %s, got %s", expected, key)
	}
}

func TestValidatePoolID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		hasError bool
	}{
		{"valid id", "aave-v3-ethereum-usdc", false},
		{"empty id", "", true},
		{"long id", string(make([]byte, 300)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidatePoolID(tt.id)
			if (len(errors) > 0) != tt.hasError {
				t.Errorf("Expected hasError=%v, got errors=%v", tt.hasError, errors)
			}
		})
	}
}

func TestValidatePeriod(t *testing.T) {
	tests := []struct {
		period   string
		hasError bool
	}{
		{"1h", false},
		{"24h", false},
		{"7d", false},
		{"30d", false},
		{"invalid", true},
		{"1w", true},
		{"", false}, // Empty is valid (will use default)
	}

	for _, tt := range tests {
		t.Run(tt.period, func(t *testing.T) {
			errors := ValidatePeriod(tt.period)
			if (len(errors) > 0) != tt.hasError {
				t.Errorf("Period %s: expected hasError=%v, got errors=%v", tt.period, tt.hasError, errors)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	err := NewAPIError(400, "BAD_REQUEST", "Invalid input")

	if err.StatusCode != 400 {
		t.Errorf("Expected status 400, got %d", err.StatusCode)
	}

	if err.Code != "BAD_REQUEST" {
		t.Errorf("Expected code BAD_REQUEST, got %s", err.Code)
	}

	if err.Error() != "Invalid input" {
		t.Errorf("Expected message 'Invalid input', got %s", err.Error())
	}

	withDetails := err.WithDetails("Field is missing")
	if withDetails.Details != "Field is missing" {
		t.Errorf("Expected details 'Field is missing', got %s", withDetails.Details)
	}
}

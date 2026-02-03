// Package utils provides common helper functions used throughout the application.
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/shopspring/decimal"
)

// GenerateID creates a random hex ID of specified length
func GenerateID(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return hex.EncodeToString([]byte(time.Now().String()))[:length]
	}
	return hex.EncodeToString(bytes)
}

// NormalizeChainName standardizes chain names
func NormalizeChainName(chain string) string {
	chain = strings.ToLower(strings.TrimSpace(chain))

	// Map common variations to standard names
	chainMap := map[string]string{
		"eth":       "ethereum",
		"mainnet":   "ethereum",
		"bnb":       "bsc",
		"binance":   "bsc",
		"matic":     "polygon",
		"poly":      "polygon",
		"arb":       "arbitrum",
		"op":        "optimism",
		"avax":      "avalanche",
		"ftm":       "fantom",
	}

	if normalized, ok := chainMap[chain]; ok {
		return normalized
	}

	return chain
}

// NormalizeProtocolName standardizes protocol names
func NormalizeProtocolName(protocol string) string {
	protocol = strings.ToLower(strings.TrimSpace(protocol))

	// Remove common suffixes
	suffixes := []string{"-v3", "-v2", "-v1", "-finance", "-protocol"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(protocol, suffix) {
			protocol = strings.TrimSuffix(protocol, suffix)
			break
		}
	}

	return protocol
}

// ParseDecimal safely parses a string to decimal
func ParseDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}

	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}

	return d
}

// FormatDecimal formats a decimal to string with specified precision
func FormatDecimal(d decimal.Decimal, precision int32) string {
	return d.StringFixed(precision)
}

// FormatPercentage formats a decimal as a percentage string
func FormatPercentage(d decimal.Decimal) string {
	return d.StringFixed(2) + "%"
}

// FormatUSD formats a decimal as USD string
func FormatUSD(d decimal.Decimal) string {
	f, _ := d.Float64()

	switch {
	case f >= 1e12:
		return "$" + decimal.NewFromFloat(f/1e12).StringFixed(2) + "T"
	case f >= 1e9:
		return "$" + decimal.NewFromFloat(f/1e9).StringFixed(2) + "B"
	case f >= 1e6:
		return "$" + decimal.NewFromFloat(f/1e6).StringFixed(2) + "M"
	case f >= 1e3:
		return "$" + decimal.NewFromFloat(f/1e3).StringFixed(2) + "K"
	default:
		return "$" + d.StringFixed(2)
	}
}

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from ends
	s = strings.Trim(s, "-")

	return s
}

// TruncateString truncates a string to max length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

// IsValidSymbol checks if a token symbol is valid
func IsValidSymbol(symbol string) bool {
	if len(symbol) == 0 || len(symbol) > 20 {
		return false
	}

	for _, r := range symbol {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}

// ContainsString checks if a slice contains a string
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UniqueStrings returns a slice with duplicate strings removed
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))

	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

// TimeAgo formats a time as a human-readable "time ago" string
func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return strings.Replace(duration.Truncate(time.Minute).String(), "0s", "", 1) + " ago"
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return strings.Replace(duration.Truncate(time.Hour).String(), "0m0s", "", 1) + " ago"
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return string(rune(days)) + " days ago"
	default:
		return t.Format("Jan 2, 2006")
	}
}

// ParseDuration parses a duration string (e.g., "3m", "1h", "24h")
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// MinDecimal returns the smaller of two decimals
func MinDecimal(a, b decimal.Decimal) decimal.Decimal {
	if a.LessThan(b) {
		return a
	}
	return b
}

// MaxDecimal returns the larger of two decimals
func MaxDecimal(a, b decimal.Decimal) decimal.Decimal {
	if a.GreaterThan(b) {
		return a
	}
	return b
}

// ClampDecimal constrains a value between min and max
func ClampDecimal(value, min, max decimal.Decimal) decimal.Decimal {
	if value.LessThan(min) {
		return min
	}
	if value.GreaterThan(max) {
		return max
	}
	return value
}

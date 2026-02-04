// Package config provides configuration management for the DeFi Yield Aggregator.
// It loads configuration from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds all application configuration
type Config struct {
	App           AppConfig
	Server        ServerConfig
	Postgres      PostgresConfig
	Redis         RedisConfig
	ElasticSearch ElasticSearchConfig
	RateLimit     RateLimitConfig
	DeFiLlama     DeFiLlamaConfig
	CoinGecko     CoinGeckoConfig
	Worker        WorkerConfig
	Scoring       ScoringConfig
	CORS          CORSConfig
	WebSocket     WebSocketConfig
}

// AppConfig holds application-level settings
type AppConfig struct {
	Env      string // development, staging, production
	Name     string
	LogLevel string
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// PostgresConfig holds PostgreSQL connection settings
type PostgresConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	Database              string
	SSLMode               string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

// DSN returns the PostgreSQL connection string
func (c PostgresConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Database +
		" sslmode=" + c.SSLMode
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// Addr returns the Redis address in host:port format
func (c RedisConfig) Addr() string {
	return c.Host + ":" + c.Port
}

// ElasticSearchConfig holds ElasticSearch connection settings
type ElasticSearchConfig struct {
	URL      string
	Username string
	Password string
}

// RateLimitConfig holds API rate limiting settings
type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

// DeFiLlamaConfig holds DeFiLlama API settings
type DeFiLlamaConfig struct {
	BaseURL       string
	RateLimit     int           // Requests per minute
	FetchInterval time.Duration // How often to fetch data
}

// CoinGeckoConfig holds CoinGecko API settings
type CoinGeckoConfig struct {
	BaseURL       string
	APIKey        string
	RateLimit     int           // Requests per minute
	FetchInterval time.Duration // How often to fetch data
}

// WorkerConfig holds background worker settings
type WorkerConfig struct {
	OpportunityDetectInterval time.Duration
	Concurrency               int
	MinTVLThreshold           float64
	MinAPYThreshold           float64
	YieldGapMinProfit         float64
	APYJumpThreshold          float64
}

// ScoringConfig holds opportunity scoring weights
type ScoringConfig struct {
	APYWeight       float64
	TVLWeight       float64
	StabilityWeight float64
	TrendWeight     float64
}

// CORSConfig holds CORS settings
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	MaxAge         int
}

// WebSocketConfig holds WebSocket settings
type WebSocketConfig struct {
	PingInterval   time.Duration
	PongTimeout    time.Duration
	MaxMessageSize int64
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("No .env file found, using environment variables")
	}

	cfg := &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			Name:     getEnv("APP_NAME", "defi-yield-aggregator"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "3000"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Postgres: PostgresConfig{
			Host:                  getEnv("POSTGRES_HOST", "localhost"),
			Port:                  getEnv("POSTGRES_PORT", "5432"),
			User:                  getEnv("POSTGRES_USER", "defi"),
			Password:             getEnv("POSTGRES_PASSWORD", "defi_secret"),
			Database:              getEnv("POSTGRES_DB", "defi_aggregator"),
			SSLMode:               getEnv("POSTGRES_SSL_MODE", "disable"),
			MaxConnections:        getInt("POSTGRES_MAX_CONNECTIONS", 25),
			MaxIdleConnections:    getInt("POSTGRES_MAX_IDLE_CONNECTIONS", 5),
			ConnectionMaxLifetime: getDuration("POSTGRES_CONNECTION_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
			PoolSize: getInt("REDIS_POOL_SIZE", 10),
		},
		ElasticSearch: ElasticSearchConfig{
			URL:      getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
			Username: getEnv("ELASTICSEARCH_USERNAME", ""),
			Password: getEnv("ELASTICSEARCH_PASSWORD", ""),
		},
		RateLimit: RateLimitConfig{
			Requests: getInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		DeFiLlama: DeFiLlamaConfig{
			BaseURL:       getEnv("DEFILLAMA_BASE_URL", "https://yields.llama.fi"),
			RateLimit:     getInt("DEFILLAMA_RATE_LIMIT", 500),
			FetchInterval: getDuration("DEFILLAMA_FETCH_INTERVAL", 3*time.Minute),
		},
		CoinGecko: CoinGeckoConfig{
			BaseURL:       getEnv("COINGECKO_BASE_URL", "https://api.coingecko.com/api/v3"),
			APIKey:        getEnv("COINGECKO_API_KEY", ""),
			RateLimit:     getInt("COINGECKO_RATE_LIMIT", 30),
			FetchInterval: getDuration("COINGECKO_FETCH_INTERVAL", 10*time.Minute),
		},
		Worker: WorkerConfig{
			OpportunityDetectInterval: getDuration("OPPORTUNITY_DETECT_INTERVAL", 5*time.Minute),
			Concurrency:               getInt("WORKER_CONCURRENCY", 5),
			MinTVLThreshold:           getFloat("MIN_TVL_THRESHOLD", 100000),
			MinAPYThreshold:           getFloat("MIN_APY_THRESHOLD", 0.1),
			YieldGapMinProfit:         getFloat("YIELD_GAP_MIN_PROFIT", 0.5),
			APYJumpThreshold:          getFloat("APY_JUMP_THRESHOLD", 50),
		},
		Scoring: ScoringConfig{
			APYWeight:       getFloat("SCORE_WEIGHT_APY", 0.35),
			TVLWeight:       getFloat("SCORE_WEIGHT_TVL", 0.25),
			StabilityWeight: getFloat("SCORE_WEIGHT_STABILITY", 0.25),
			TrendWeight:     getFloat("SCORE_WEIGHT_TREND", 0.15),
		},
		CORS: CORSConfig{
			AllowedOrigins: getStringSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getStringSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getStringSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
			MaxAge:         getInt("CORS_MAX_AGE", 86400),
		},
		WebSocket: WebSocketConfig{
			PingInterval:   getDuration("WS_PING_INTERVAL", 30*time.Second),
			PongTimeout:    getDuration("WS_PONG_TIMEOUT", 60*time.Second),
			MaxMessageSize: int64(getInt("WS_MAX_MESSAGE_SIZE", 65536)), // 64KB for pool updates
		},
	}

	return cfg, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// Helper functions for reading environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}

-- =============================================================================
-- DeFi Yield Aggregator - Database Schema
-- Migration: 001_create_tables
-- =============================================================================
-- This migration creates the core database tables for the DeFi Yield Aggregator.
-- It uses TimescaleDB for efficient time-series data storage.

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- =============================================================================
-- Pools Table
-- =============================================================================
-- Stores current state of all DeFi yield pools
CREATE TABLE IF NOT EXISTS pools (
    id VARCHAR(255) PRIMARY KEY,               -- Unique pool identifier (from DeFiLlama)
    chain VARCHAR(50) NOT NULL,                -- Blockchain network (ethereum, bsc, polygon, etc.)
    protocol VARCHAR(100) NOT NULL,            -- Protocol name (aave-v3, compound, curve, etc.)
    symbol VARCHAR(100) NOT NULL,              -- Pool symbol (USDC, ETH-USDC, etc.)

    -- Financial metrics
    tvl DECIMAL(24, 2) NOT NULL DEFAULT 0,     -- Total Value Locked in USD
    apy DECIMAL(12, 6) NOT NULL DEFAULT 0,     -- Current APY (percentage)
    apy_base DECIMAL(12, 6) DEFAULT 0,         -- Base APY (lending/trading fees)
    apy_reward DECIMAL(12, 6) DEFAULT 0,       -- Reward APY (token incentives)

    -- Tokens
    reward_tokens TEXT[] DEFAULT '{}',          -- Reward token symbols
    underlying_tokens TEXT[] DEFAULT '{}',      -- Underlying asset symbols
    pool_meta TEXT,                             -- Additional metadata

    -- Risk metrics
    il_7d DECIMAL(12, 6) DEFAULT 0,            -- 7-day impermanent loss
    apy_mean_30d DECIMAL(12, 6) DEFAULT 0,     -- 30-day average APY

    -- Volume metrics
    volume_usd_1d DECIMAL(24, 2) DEFAULT 0,    -- 24h trading volume
    volume_usd_7d DECIMAL(24, 2) DEFAULT 0,    -- 7-day trading volume

    -- Calculated metrics
    score DECIMAL(6, 2) DEFAULT 0,             -- Risk-adjusted score (0-100)
    apy_change_1h DECIMAL(12, 6) DEFAULT 0,    -- APY change in last hour
    apy_change_24h DECIMAL(12, 6) DEFAULT 0,   -- APY change in last 24 hours
    apy_change_7d DECIMAL(12, 6) DEFAULT 0,    -- APY change in last 7 days

    -- Categories
    stablecoin BOOLEAN DEFAULT FALSE,          -- Is stablecoin pool
    exposure VARCHAR(20) DEFAULT 'single',     -- single, multi

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_pools_chain ON pools(chain);
CREATE INDEX IF NOT EXISTS idx_pools_protocol ON pools(protocol);
CREATE INDEX IF NOT EXISTS idx_pools_tvl ON pools(tvl DESC);
CREATE INDEX IF NOT EXISTS idx_pools_apy ON pools(apy DESC);
CREATE INDEX IF NOT EXISTS idx_pools_score ON pools(score DESC);
CREATE INDEX IF NOT EXISTS idx_pools_stablecoin ON pools(stablecoin);
CREATE INDEX IF NOT EXISTS idx_pools_chain_protocol ON pools(chain, protocol);
CREATE INDEX IF NOT EXISTS idx_pools_updated_at ON pools(updated_at);

-- =============================================================================
-- Historical APY Table (TimescaleDB Hypertable)
-- =============================================================================
-- Stores time-series APY and TVL data for historical analysis
CREATE TABLE IF NOT EXISTS historical_apy (
    pool_id VARCHAR(255) NOT NULL,             -- Reference to pools.id
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    apy DECIMAL(12, 6) NOT NULL,
    tvl DECIMAL(24, 2) NOT NULL,
    apy_base DECIMAL(12, 6) DEFAULT 0,
    apy_reward DECIMAL(12, 6) DEFAULT 0,

    PRIMARY KEY (pool_id, timestamp)
);

-- Convert to TimescaleDB hypertable for efficient time-series queries
-- Partition by time (7-day chunks)
SELECT create_hypertable(
    'historical_apy',
    'timestamp',
    chunk_time_interval => INTERVAL '7 days',
    if_not_exists => TRUE
);

-- Create index for fast lookups by pool
CREATE INDEX IF NOT EXISTS idx_historical_apy_pool_time
    ON historical_apy(pool_id, timestamp DESC);

-- Enable compression for older data (after 30 days)
ALTER TABLE historical_apy SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'pool_id'
);

-- Add compression policy (compress chunks older than 30 days)
SELECT add_compression_policy('historical_apy', INTERVAL '30 days', if_not_exists => TRUE);

-- Add retention policy (delete data older than 1 year)
SELECT add_retention_policy('historical_apy', INTERVAL '365 days', if_not_exists => TRUE);

-- =============================================================================
-- Opportunities Table
-- =============================================================================
-- Stores detected yield farming opportunities
CREATE TABLE IF NOT EXISTS opportunities (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(20) NOT NULL,                 -- yield-gap, trending, high-score
    title VARCHAR(255) NOT NULL,
    description TEXT,

    -- Pool references
    source_pool_id VARCHAR(255),               -- For yield-gap: source pool
    target_pool_id VARCHAR(255),               -- For yield-gap: target pool
    pool_id VARCHAR(255),                      -- For trending/high-score: single pool

    -- Metrics
    asset VARCHAR(50),                         -- Primary asset (USDC, ETH, etc.)
    chain VARCHAR(50),
    apy_difference DECIMAL(12, 6) DEFAULT 0,   -- For yield-gap
    apy_growth DECIMAL(12, 6) DEFAULT 0,       -- For trending
    current_apy DECIMAL(12, 6) DEFAULT 0,
    potential_profit DECIMAL(12, 6) DEFAULT 0, -- Estimated profit percentage
    tvl DECIMAL(24, 2) DEFAULT 0,

    -- Risk assessment
    risk_level VARCHAR(10) DEFAULT 'medium',   -- low, medium, high
    score DECIMAL(6, 2) DEFAULT 0,             -- Opportunity score (0-100)

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    detected_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_seen_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for opportunity queries
CREATE INDEX IF NOT EXISTS idx_opportunities_type ON opportunities(type);
CREATE INDEX IF NOT EXISTS idx_opportunities_chain ON opportunities(chain);
CREATE INDEX IF NOT EXISTS idx_opportunities_asset ON opportunities(asset);
CREATE INDEX IF NOT EXISTS idx_opportunities_risk_level ON opportunities(risk_level);
CREATE INDEX IF NOT EXISTS idx_opportunities_is_active ON opportunities(is_active);
CREATE INDEX IF NOT EXISTS idx_opportunities_score ON opportunities(score DESC);
CREATE INDEX IF NOT EXISTS idx_opportunities_detected_at ON opportunities(detected_at DESC);
CREATE INDEX IF NOT EXISTS idx_opportunities_expires_at ON opportunities(expires_at);

-- =============================================================================
-- Token Prices Table
-- =============================================================================
-- Caches token prices from CoinGecko for profitability calculations
CREATE TABLE IF NOT EXISTS token_prices (
    token_id VARCHAR(100) PRIMARY KEY,         -- CoinGecko token ID
    symbol VARCHAR(20) NOT NULL,
    price_usd DECIMAL(24, 8) NOT NULL,
    market_cap DECIMAL(24, 2),
    volume_24h DECIMAL(24, 2),
    price_change_24h DECIMAL(12, 6),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_token_prices_symbol ON token_prices(symbol);

-- =============================================================================
-- Functions and Triggers
-- =============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for pools table
DROP TRIGGER IF EXISTS update_pools_updated_at ON pools;
CREATE TRIGGER update_pools_updated_at
    BEFORE UPDATE ON pools
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for opportunities table
DROP TRIGGER IF EXISTS update_opportunities_updated_at ON opportunities;
CREATE TRIGGER update_opportunities_updated_at
    BEFORE UPDATE ON opportunities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- Continuous Aggregates (TimescaleDB)
-- =============================================================================
-- Pre-computed aggregations for faster queries

-- Hourly APY aggregates
CREATE MATERIALIZED VIEW IF NOT EXISTS hourly_apy_stats
WITH (timescaledb.continuous) AS
SELECT
    pool_id,
    time_bucket('1 hour', timestamp) AS hour,
    AVG(apy) AS avg_apy,
    MAX(apy) AS max_apy,
    MIN(apy) AS min_apy,
    AVG(tvl) AS avg_tvl,
    COUNT(*) AS data_points
FROM historical_apy
GROUP BY pool_id, time_bucket('1 hour', timestamp)
WITH NO DATA;

-- Refresh policy for hourly stats
SELECT add_continuous_aggregate_policy('hourly_apy_stats',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

-- Daily APY aggregates
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_apy_stats
WITH (timescaledb.continuous) AS
SELECT
    pool_id,
    time_bucket('1 day', timestamp) AS day,
    AVG(apy) AS avg_apy,
    MAX(apy) AS max_apy,
    MIN(apy) AS min_apy,
    AVG(tvl) AS avg_tvl,
    STDDEV(apy) AS apy_stddev,
    COUNT(*) AS data_points
FROM historical_apy
GROUP BY pool_id, time_bucket('1 day', timestamp)
WITH NO DATA;

-- Refresh policy for daily stats
SELECT add_continuous_aggregate_policy('daily_apy_stats',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

-- =============================================================================
-- Grant Permissions (for production)
-- =============================================================================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO defi_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO defi_app;

COMMENT ON TABLE pools IS 'Current state of all DeFi yield pools from DeFiLlama';
COMMENT ON TABLE historical_apy IS 'Time-series APY and TVL data for historical analysis';
COMMENT ON TABLE opportunities IS 'Detected yield farming opportunities';
COMMENT ON TABLE token_prices IS 'Cached token prices from CoinGecko';

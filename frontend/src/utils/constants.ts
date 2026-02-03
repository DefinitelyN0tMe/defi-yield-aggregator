/**
 * Application constants
 */

// API Configuration
export const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1';
export const USE_MOCK_DATA = import.meta.env.VITE_USE_MOCK_DATA !== 'false';

// Query Configuration
export const QUERY_STALE_TIME = 30000; // 30 seconds
export const QUERY_CACHE_TIME = 5 * 60 * 1000; // 5 minutes
export const QUERY_RETRY_COUNT = 2;

// Refetch Intervals
export const STATS_REFETCH_INTERVAL = 30000; // 30 seconds
export const POOLS_REFETCH_INTERVAL = 60000; // 1 minute

// Pagination
export const DEFAULT_PAGE_SIZE = 25;
export const OPPORTUNITIES_PAGE_SIZE = 20;
export const TOP_POOLS_LIMIT = 6;
export const TRENDING_LIMIT = 5;

// TVL Filter Thresholds
export const TVL_THRESHOLDS = {
  SMALL: 100_000,
  MEDIUM: 1_000_000,
  LARGE: 10_000_000,
  WHALE: 100_000_000,
} as const;

// Chart Configuration
export const CHART_COLORS = {
  GREEN: '#22c55e',
  RED: '#ef4444',
  PRIMARY: '#6366f1',
  GRAY: '#374151',
} as const;

// Score Thresholds
export const SCORE_THRESHOLDS = {
  HIGH: 80,
  MEDIUM: 60,
} as const;

// WebSocket Configuration
export const WS_RECONNECT_DELAY = 5000; // 5 seconds
export const WS_BASE_URL = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080';

// Protocols list for detection
export const KNOWN_PROTOCOLS = [
  'aave',
  'compound',
  'uniswap',
  'curve',
  'lido',
  'gmx',
  'pendle',
  'velodrome',
  'yearn',
  'beefy',
  'convex',
  'rocket-pool',
  'maker',
  'sushiswap',
] as const;

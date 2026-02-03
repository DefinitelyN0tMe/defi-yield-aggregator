// Pool types
export interface Pool {
  id: string;
  chain: string;
  protocol: string;
  symbol: string;
  tvl: number;
  apy: number;
  apyBase: number;
  apyReward: number;
  rewardTokens: string[];
  underlyingTokens: string[];
  poolMeta: string;
  il7d: number;
  apyMean30d: number;
  volumeUsd1d: number;
  volumeUsd7d: number;
  score: number;
  apyChange1h: number;
  apyChange24h: number;
  apyChange7d: number;
  stablecoin: boolean;
  exposure: string;
  createdAt: string;
  updatedAt: string;
}

export interface PoolListResponse {
  data: Pool[];
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

export interface PoolFilter {
  chain?: string;
  protocol?: string;
  symbol?: string;
  search?: string; // Search by symbol, protocol, or chain
  minApy?: number;
  maxApy?: number;
  minTvl?: number;
  maxTvl?: number;
  minScore?: number;
  stablecoin?: boolean;
  sortBy?: 'apy' | 'tvl' | 'score';
  sortOrder?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

// Historical APY
export interface HistoricalAPY {
  poolId: string;
  timestamp: string;
  apy: number;
  tvl: number;
  apyBase: number;
  apyReward: number;
}

export interface PoolHistoryResponse {
  poolId: string;
  period: string;
  dataPoints: HistoricalAPY[];
}

// Opportunity types
export type OpportunityType = 'yield-gap' | 'trending' | 'high-score';
export type RiskLevel = 'low' | 'medium' | 'high';

export interface Opportunity {
  id: string;
  type: OpportunityType;
  title: string;
  description: string;
  sourcePoolId?: string;
  targetPoolId?: string;
  poolId?: string;
  asset: string;
  chain: string;
  apyDifference: number;
  apyGrowth: number;
  currentApy: number;
  potentialProfit: number;
  tvl: number;
  riskLevel: RiskLevel;
  score: number;
  isActive: boolean;
  detectedAt: string;
  lastSeenAt: string;
  expiresAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface OpportunityListResponse {
  data: Opportunity[];
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

export interface OpportunityFilter {
  type?: OpportunityType;
  riskLevel?: RiskLevel;
  chain?: string;
  asset?: string;
  minProfit?: number;
  minScore?: number;
  activeOnly?: boolean;
  sortBy?: 'score' | 'profit' | 'apy' | 'detectedAt';
  sortOrder?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

// Trending pools
export interface TrendingPool {
  pool: Pool;
  apyGrowth1h: number;
  apyGrowth24h: number;
  apyGrowth7d: number;
  trendScore: number;
}

// Chain and Protocol types
export interface Chain {
  name: string;
  displayName: string;
  poolCount: number;
  totalTvl: number;
  averageApy: number;
  maxApy: number;
  topProtocols?: string[];
}

export interface Protocol {
  name: string;
  displayName: string;
  category?: string;
  chains: string[];
  poolCount: number;
  totalTvl: number;
  averageApy: number;
  maxApy: number;
}

// Platform stats
export interface APYDistribution {
  range0to1: number;
  range1to5: number;
  range5to10: number;
  range10to25: number;
  range25to50: number;
  range50to100: number;
  range100plus: number;
}

export interface PlatformStats {
  totalPools: number;
  totalTvl: number;
  averageApy: number;
  medianApy: number;
  maxApy: number;
  totalChains: number;
  totalProtocols: number;
  activeOpportunities: number;
  lastUpdated: string;
  tvlByChain: Record<string, number>;
  poolsByChain: Record<string, number>;
  apyDistribution: APYDistribution;
}

// Health check
export interface ServiceHealth {
  status: 'up' | 'down';
  latency: string;
  message?: string;
}

export interface HealthCheck {
  status: 'healthy' | 'degraded' | 'unhealthy';
  version: string;
  uptime: string;
  timestamp: string;
  services: {
    postgresql: ServiceHealth;
    redis: ServiceHealth;
    elasticsearch: ServiceHealth;
  };
}

// WebSocket message types
export interface WSMessage {
  type: 'pool_update' | 'opportunity_alert' | 'ping' | 'pong';
  timestamp: string;
  data?: Pool | Opportunity;
}

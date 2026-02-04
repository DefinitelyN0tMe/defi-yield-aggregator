import type {
  Pool,
  PoolListResponse,
  PoolFilter,
  PoolHistoryResponse,
  Opportunity,
  OpportunityListResponse,
  OpportunityFilter,
  TrendingPool,
  Chain,
  Protocol,
  PlatformStats,
  HealthCheck,
  HistoricalAPY,
} from '../types';
import {
  mockPools,
  mockOpportunities,
  mockStats,
  mockChains,
  mockTrendingPools,
  generateMockHistory,
} from './mockData';
import { API_BASE, USE_MOCK_DATA } from '../utils/constants';

// Helper to convert string to number safely
function toNumber(value: unknown): number {
  if (typeof value === 'number') return value;
  if (typeof value === 'string') {
    const num = parseFloat(value);
    return isNaN(num) ? 0 : num;
  }
  return 0;
}

// Transform raw API pool data to proper typed Pool
function transformPool(raw: Record<string, unknown>): Pool {
  return {
    id: String(raw.id || ''),
    chain: String(raw.chain || ''),
    protocol: String(raw.protocol || ''),
    symbol: String(raw.symbol || ''),
    tvl: toNumber(raw.tvl),
    apy: toNumber(raw.apy),
    apyBase: toNumber(raw.apyBase || raw.apy_base),
    apyReward: toNumber(raw.apyReward || raw.apy_reward),
    rewardTokens: Array.isArray(raw.rewardTokens || raw.reward_tokens)
      ? (raw.rewardTokens || raw.reward_tokens) as string[]
      : [],
    underlyingTokens: Array.isArray(raw.underlyingTokens || raw.underlying_tokens)
      ? (raw.underlyingTokens || raw.underlying_tokens) as string[]
      : [],
    poolMeta: String(raw.poolMeta || raw.pool_meta || ''),
    il7d: toNumber(raw.il7d || raw.il_7d),
    apyMean30d: toNumber(raw.apyMean30d || raw.apy_mean_30d),
    volumeUsd1d: toNumber(raw.volumeUsd1d || raw.volume_usd_1d),
    volumeUsd7d: toNumber(raw.volumeUsd7d || raw.volume_usd_7d),
    score: toNumber(raw.score),
    apyChange1h: toNumber(raw.apyChange1h || raw.apy_change_1h),
    apyChange24h: toNumber(raw.apyChange24h || raw.apy_change_24h),
    apyChange7d: toNumber(raw.apyChange7d || raw.apy_change_7d),
    stablecoin: Boolean(raw.stablecoin),
    exposure: String(raw.exposure || ''),
    createdAt: String(raw.createdAt || raw.created_at || ''),
    updatedAt: String(raw.updatedAt || raw.updated_at || ''),
  };
}

// Transform raw API chain data
function transformChain(raw: Record<string, unknown>): Chain {
  return {
    name: String(raw.name || ''),
    displayName: String(raw.displayName || raw.display_name || raw.name || ''),
    poolCount: toNumber(raw.poolCount || raw.pool_count),
    totalTvl: toNumber(raw.totalTvl || raw.total_tvl),
    averageApy: toNumber(raw.averageApy || raw.average_apy),
    maxApy: toNumber(raw.maxApy || raw.max_apy),
    topProtocols: Array.isArray(raw.topProtocols || raw.top_protocols)
      ? (raw.topProtocols || raw.top_protocols) as string[]
      : undefined,
  };
}

// Transform raw API protocol data
function transformProtocol(raw: Record<string, unknown>): Protocol {
  return {
    name: String(raw.name || ''),
    displayName: String(raw.displayName || raw.display_name || raw.name || ''),
    category: raw.category ? String(raw.category) : undefined,
    chains: Array.isArray(raw.chains) ? raw.chains as string[] : [],
    poolCount: toNumber(raw.poolCount || raw.pool_count),
    totalTvl: toNumber(raw.totalTvl || raw.total_tvl),
    averageApy: toNumber(raw.averageApy || raw.average_apy),
    maxApy: toNumber(raw.maxApy || raw.max_apy),
  };
}

// Transform raw API opportunity data
function transformOpportunity(raw: Record<string, unknown>): Opportunity {
  return {
    id: String(raw.id || ''),
    type: (raw.type || 'high-score') as Opportunity['type'],
    title: String(raw.title || ''),
    description: String(raw.description || ''),
    sourcePoolId: raw.sourcePoolId || raw.source_pool_id ? String(raw.sourcePoolId || raw.source_pool_id) : undefined,
    targetPoolId: raw.targetPoolId || raw.target_pool_id ? String(raw.targetPoolId || raw.target_pool_id) : undefined,
    poolId: raw.poolId || raw.pool_id ? String(raw.poolId || raw.pool_id) : undefined,
    asset: String(raw.asset || ''),
    chain: String(raw.chain || ''),
    apyDifference: toNumber(raw.apyDifference || raw.apy_difference),
    apyGrowth: toNumber(raw.apyGrowth || raw.apy_growth),
    currentApy: toNumber(raw.currentApy || raw.current_apy),
    potentialProfit: toNumber(raw.potentialProfit || raw.potential_profit),
    tvl: toNumber(raw.tvl),
    riskLevel: (raw.riskLevel || raw.risk_level || 'medium') as Opportunity['riskLevel'],
    score: toNumber(raw.score),
    isActive: Boolean(raw.isActive ?? raw.is_active ?? true),
    detectedAt: String(raw.detectedAt || raw.detected_at || ''),
    lastSeenAt: String(raw.lastSeenAt || raw.last_seen_at || ''),
    expiresAt: String(raw.expiresAt || raw.expires_at || ''),
    createdAt: String(raw.createdAt || raw.created_at || ''),
    updatedAt: String(raw.updatedAt || raw.updated_at || ''),
  };
}

// Transform raw API stats data
function transformStats(raw: Record<string, unknown>): PlatformStats {
  const tvlByChain: Record<string, number> = {};
  const rawTvlByChain = (raw.tvlByChain || raw.tvl_by_chain || {}) as Record<string, unknown>;
  for (const [key, value] of Object.entries(rawTvlByChain)) {
    tvlByChain[key] = toNumber(value);
  }

  const poolsByChain: Record<string, number> = {};
  const rawPoolsByChain = (raw.poolsByChain || raw.pools_by_chain || {}) as Record<string, unknown>;
  for (const [key, value] of Object.entries(rawPoolsByChain)) {
    poolsByChain[key] = toNumber(value);
  }

  const rawApyDist = (raw.apyDistribution || raw.apy_distribution || {}) as Record<string, unknown>;

  return {
    totalPools: toNumber(raw.totalPools || raw.total_pools),
    totalTvl: toNumber(raw.totalTvl || raw.total_tvl),
    averageApy: toNumber(raw.averageApy || raw.average_apy),
    medianApy: toNumber(raw.medianApy || raw.median_apy),
    maxApy: toNumber(raw.maxApy || raw.max_apy),
    totalChains: toNumber(raw.totalChains || raw.total_chains),
    totalProtocols: toNumber(raw.totalProtocols || raw.total_protocols),
    activeOpportunities: toNumber(raw.activeOpportunities || raw.active_opportunities),
    lastUpdated: String(raw.lastUpdated || raw.last_updated || new Date().toISOString()),
    tvlByChain,
    poolsByChain,
    apyDistribution: {
      range0to1: toNumber(rawApyDist.range0to1 || rawApyDist['0-1'] || rawApyDist.range_0_to_1),
      range1to5: toNumber(rawApyDist.range1to5 || rawApyDist['1-5'] || rawApyDist.range_1_to_5),
      range5to10: toNumber(rawApyDist.range5to10 || rawApyDist['5-10'] || rawApyDist.range_5_to_10),
      range10to25: toNumber(rawApyDist.range10to25 || rawApyDist['10-25'] || rawApyDist.range_10_to_25),
      range25to50: toNumber(rawApyDist.range25to50 || rawApyDist['25-50'] || rawApyDist.range_25_to_50),
      range50to100: toNumber(rawApyDist.range50to100 || rawApyDist['50-100'] || rawApyDist.range_50_to_100),
      range100plus: toNumber(rawApyDist.range100plus || rawApyDist['100+'] || rawApyDist.range_100_plus),
    },
  };
}

// Transform trending pool data
function transformTrendingPool(raw: Record<string, unknown>): TrendingPool {
  return {
    pool: transformPool((raw.pool || {}) as Record<string, unknown>),
    apyGrowth1h: toNumber(raw.apyGrowth1h || raw.apy_growth_1h),
    apyGrowth24h: toNumber(raw.apyGrowth24h || raw.apy_growth_24h),
    apyGrowth7d: toNumber(raw.apyGrowth7d || raw.apy_growth_7d),
    trendScore: toNumber(raw.trendScore || raw.trend_score),
  };
}

// Transform historical APY data
function transformHistoricalAPY(raw: Record<string, unknown>): HistoricalAPY {
  return {
    poolId: String(raw.poolId || raw.pool_id || ''),
    timestamp: String(raw.timestamp || ''),
    apy: toNumber(raw.apy),
    tvl: toNumber(raw.tvl),
    apyBase: toNumber(raw.apyBase || raw.apy_base),
    apyReward: toNumber(raw.apyReward || raw.apy_reward),
  };
}

// Helper to build query string from filter object
function buildQueryString(params: Record<string, unknown>): string {
  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.append(key, String(value));
    }
  });

  const query = searchParams.toString();
  return query ? `?${query}` : '';
}

// Generic fetch wrapper with error handling
async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'An error occurred' }));
    throw new Error(error.error?.message || error.message || `HTTP ${response.status}`);
  }

  try {
    return await response.json();
  } catch (err) {
    console.error('Failed to parse API response:', err);
    throw new Error('Invalid response from server');
  }
}

// Helper to filter and sort pools
function filterPools(pools: Pool[], filter?: PoolFilter): Pool[] {
  let result = [...pools];

  // Search filter - searches across symbol, protocol, and chain
  if (filter?.search) {
    const searchLower = filter.search.toLowerCase();
    result = result.filter(p =>
      p.symbol.toLowerCase().includes(searchLower) ||
      p.protocol.toLowerCase().includes(searchLower) ||
      p.chain.toLowerCase().includes(searchLower) ||
      p.poolMeta?.toLowerCase().includes(searchLower)
    );
  }

  if (filter?.chain) {
    result = result.filter(p => p.chain === filter.chain);
  }
  if (filter?.protocol) {
    result = result.filter(p => p.protocol === filter.protocol);
  }
  if (filter?.minApy !== undefined) {
    result = result.filter(p => p.apy >= filter.minApy!);
  }
  if (filter?.maxApy !== undefined) {
    result = result.filter(p => p.apy <= filter.maxApy!);
  }
  if (filter?.minTvl !== undefined) {
    result = result.filter(p => p.tvl >= filter.minTvl!);
  }
  if (filter?.stablecoin !== undefined) {
    result = result.filter(p => p.stablecoin === filter.stablecoin);
  }

  // Sort
  const sortBy = filter?.sortBy || 'score';
  const sortOrder = filter?.sortOrder || 'desc';
  result.sort((a, b) => {
    const aVal = a[sortBy as keyof Pool] as number;
    const bVal = b[sortBy as keyof Pool] as number;
    return sortOrder === 'desc' ? bVal - aVal : aVal - bVal;
  });

  return result;
}

// Helper to filter opportunities
function filterOpportunities(opps: Opportunity[], filter?: OpportunityFilter): Opportunity[] {
  let result = [...opps];

  if (filter?.type) {
    result = result.filter(o => o.type === filter.type);
  }
  if (filter?.riskLevel) {
    result = result.filter(o => o.riskLevel === filter.riskLevel);
  }
  if (filter?.chain) {
    result = result.filter(o => o.chain === filter.chain);
  }
  if (filter?.activeOnly) {
    result = result.filter(o => o.isActive);
  }

  // Sort
  const sortBy = filter?.sortBy || 'score';
  const sortOrder = filter?.sortOrder || 'desc';
  result.sort((a, b) => {
    let aVal: number, bVal: number;
    switch (sortBy) {
      case 'profit':
        aVal = a.potentialProfit;
        bVal = b.potentialProfit;
        break;
      case 'apy':
        aVal = a.currentApy;
        bVal = b.currentApy;
        break;
      case 'detectedAt':
        aVal = new Date(a.detectedAt).getTime();
        bVal = new Date(b.detectedAt).getTime();
        break;
      default:
        aVal = a.score;
        bVal = b.score;
    }
    return sortOrder === 'desc' ? bVal - aVal : aVal - bVal;
  });

  return result;
}

// Pool API
export const poolsApi = {
  list: async (filter?: PoolFilter): Promise<PoolListResponse> => {
    if (USE_MOCK_DATA) {
      const filtered = filterPools(mockPools, filter);
      const offset = filter?.offset || 0;
      const limit = filter?.limit || 25;
      const paginated = filtered.slice(offset, offset + limit);
      return {
        data: paginated,
        total: filtered.length,
        limit,
        offset,
        hasMore: offset + limit < filtered.length,
      };
    }
    const query = buildQueryString((filter || {}) as Record<string, unknown>);
    const raw = await fetchApi<{ data: Record<string, unknown>[]; total: unknown; limit: unknown; offset: unknown; hasMore: unknown }>(`/pools${query}`);
    return {
      data: (raw.data || []).map(transformPool),
      total: toNumber(raw.total),
      limit: toNumber(raw.limit),
      offset: toNumber(raw.offset),
      hasMore: Boolean(raw.hasMore),
    };
  },

  get: async (id: string): Promise<Pool> => {
    if (USE_MOCK_DATA) {
      const pool = mockPools.find(p => p.id === id);
      if (!pool) throw new Error('Pool not found');
      return pool;
    }
    const raw = await fetchApi<Record<string, unknown>>(`/pools/${encodeURIComponent(id)}`);
    return transformPool(raw);
  },

  getHistory: async (id: string, period: '1h' | '24h' | '7d' | '30d' = '7d'): Promise<PoolHistoryResponse> => {
    if (USE_MOCK_DATA) {
      return generateMockHistory(id, period);
    }
    const raw = await fetchApi<{ poolId?: string; pool_id?: string; period?: string; dataPoints?: Record<string, unknown>[]; data_points?: Record<string, unknown>[] }>(`/pools/${encodeURIComponent(id)}/history?period=${period}`);
    return {
      poolId: String(raw.poolId || raw.pool_id || id),
      period: String(raw.period || period),
      dataPoints: (raw.dataPoints || raw.data_points || []).map(transformHistoricalAPY),
    };
  },
};

// Opportunities API
export const opportunitiesApi = {
  list: async (filter?: OpportunityFilter): Promise<OpportunityListResponse> => {
    if (USE_MOCK_DATA) {
      const filtered = filterOpportunities(mockOpportunities, filter);
      const offset = filter?.offset || 0;
      const limit = filter?.limit || 20;
      const paginated = filtered.slice(offset, offset + limit);
      return {
        data: paginated,
        total: filtered.length,
        limit,
        offset,
        hasMore: offset + limit < filtered.length,
      };
    }
    const query = buildQueryString((filter || {}) as Record<string, unknown>);
    const raw = await fetchApi<{ data: Record<string, unknown>[]; total: unknown; limit: unknown; offset: unknown; hasMore: unknown }>(`/opportunities${query}`);
    return {
      data: (raw.data || []).map(transformOpportunity),
      total: toNumber(raw.total),
      limit: toNumber(raw.limit),
      offset: toNumber(raw.offset),
      hasMore: Boolean(raw.hasMore),
    };
  },

  getTrending: async (params?: { chain?: string; minGrowth?: number; limit?: number }): Promise<{ data: TrendingPool[] }> => {
    if (USE_MOCK_DATA) {
      let result = [...mockTrendingPools];
      if (params?.chain) {
        result = result.filter(t => t.pool.chain === params.chain);
      }
      if (params?.limit) {
        result = result.slice(0, params.limit);
      }
      return { data: result };
    }
    const query = buildQueryString(params || {});
    const raw = await fetchApi<{ data: Record<string, unknown>[] }>(`/opportunities/trending${query}`);
    return {
      data: (raw.data || []).map(transformTrendingPool),
    };
  },
};

// Chains API
export const chainsApi = {
  list: async (): Promise<{ data: Chain[]; total: number }> => {
    if (USE_MOCK_DATA) {
      return { data: mockChains, total: mockChains.length };
    }
    const raw = await fetchApi<{ data: Record<string, unknown>[]; total: unknown }>('/chains');
    return {
      data: (raw.data || []).map(transformChain),
      total: toNumber(raw.total),
    };
  },
};

// Protocols API
export const protocolsApi = {
  list: async (params?: { chain?: string; sortBy?: string; limit?: number }): Promise<{
    data: Protocol[];
    total: number;
    hasMore: boolean;
  }> => {
    if (USE_MOCK_DATA) {
      // Return empty for now
      return { data: [], total: 0, hasMore: false };
    }
    const query = buildQueryString(params || {});
    const raw = await fetchApi<{ data: Record<string, unknown>[]; total: unknown; hasMore: unknown }>(`/protocols${query}`);
    return {
      data: (raw.data || []).map(transformProtocol),
      total: toNumber(raw.total),
      hasMore: Boolean(raw.hasMore),
    };
  },
};

// Stats API
export const statsApi = {
  get: async (): Promise<PlatformStats> => {
    if (USE_MOCK_DATA) {
      return mockStats;
    }
    const raw = await fetchApi<Record<string, unknown>>('/stats');
    return transformStats(raw);
  },
};

// Health API
export const healthApi = {
  check: async (): Promise<HealthCheck> => {
    if (USE_MOCK_DATA) {
      return {
        status: 'healthy',
        version: '1.0.0',
        uptime: '24h 30m',
        timestamp: new Date().toISOString(),
        services: {
          postgresql: { status: 'up', latency: '2ms' },
          redis: { status: 'up', latency: '1ms' },
          elasticsearch: { status: 'up', latency: '5ms' },
        },
      };
    }
    return fetchApi<HealthCheck>('/health');
  },
};

// Export all APIs
export const api = {
  pools: poolsApi,
  opportunities: opportunitiesApi,
  chains: chainsApi,
  protocols: protocolsApi,
  stats: statsApi,
  health: healthApi,
};

export default api;

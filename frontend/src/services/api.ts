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

  return response.json();
}

// Helper to filter and sort pools
function filterPools(pools: Pool[], filter?: PoolFilter): Pool[] {
  let result = [...pools];

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
    return fetchApi<PoolListResponse>(`/pools${query}`);
  },

  get: async (id: string): Promise<Pool> => {
    if (USE_MOCK_DATA) {
      const pool = mockPools.find(p => p.id === id);
      if (!pool) throw new Error('Pool not found');
      return pool;
    }
    return fetchApi<Pool>(`/pools/${encodeURIComponent(id)}`);
  },

  getHistory: async (id: string, period: '1h' | '24h' | '7d' | '30d' = '7d'): Promise<PoolHistoryResponse> => {
    if (USE_MOCK_DATA) {
      return generateMockHistory(id, period);
    }
    return fetchApi<PoolHistoryResponse>(`/pools/${encodeURIComponent(id)}/history?period=${period}`);
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
    return fetchApi<OpportunityListResponse>(`/opportunities${query}`);
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
    return fetchApi<{ data: TrendingPool[] }>(`/opportunities/trending${query}`);
  },
};

// Chains API
export const chainsApi = {
  list: async (): Promise<{ data: Chain[]; total: number }> => {
    if (USE_MOCK_DATA) {
      return { data: mockChains, total: mockChains.length };
    }
    return fetchApi<{ data: Chain[]; total: number }>('/chains');
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
    return fetchApi<{ data: Protocol[]; total: number; hasMore: boolean }>(`/protocols${query}`);
  },
};

// Stats API
export const statsApi = {
  get: async (): Promise<PlatformStats> => {
    if (USE_MOCK_DATA) {
      return mockStats;
    }
    return fetchApi<PlatformStats>('/stats');
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

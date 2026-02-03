import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { statsApi, poolsApi, opportunitiesApi } from '../services/api';
import { StatsCard, PoolCard, OpportunityCard, TVLIcon, PoolsIcon, APYIcon, OpportunitiesIcon } from '../components';
import { usePoolUpdates } from '../hooks/useWebSocket';
import { formatNumber } from '../utils/format';
import { STATS_REFETCH_INTERVAL, POOLS_REFETCH_INTERVAL, TOP_POOLS_LIMIT } from '../utils/constants';
import type { Opportunity } from '../types';

export function HomePage() {
  const navigate = useNavigate();

  // Fetch platform stats
  const { data: stats } = useQuery({
    queryKey: ['stats'],
    queryFn: statsApi.get,
    refetchInterval: STATS_REFETCH_INTERVAL,
  });

  // Fetch top pools
  const { data: poolsData, isLoading: poolsLoading } = useQuery({
    queryKey: ['pools', 'top'],
    queryFn: () => poolsApi.list({ sortBy: 'score', sortOrder: 'desc', limit: TOP_POOLS_LIMIT }),
    refetchInterval: POOLS_REFETCH_INTERVAL,
  });

  // Fetch active opportunities
  const { data: oppsData, isLoading: oppsLoading } = useQuery({
    queryKey: ['opportunities', 'active'],
    queryFn: () => opportunitiesApi.list({ activeOnly: true, limit: 4 }),
    refetchInterval: POOLS_REFETCH_INTERVAL,
  });

  // Subscribe to real-time pool updates (callback can be used for notifications)
  usePoolUpdates(() => {
    // Real-time pool updates received - could trigger UI refresh
  });

  // Handle opportunity click - navigate to opportunity details
  const handleOpportunityClick = (opportunity: Opportunity) => {
    navigate(`/opportunities/${encodeURIComponent(opportunity.id)}`);
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Dashboard</h1>
        <p className="mt-2 text-gray-400">
          Track and analyze yield farming opportunities across DeFi protocols
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Total Value Locked"
          value={formatNumber(stats?.totalTvl)}
          subtitle={`Across ${stats?.totalChains || 0} chains`}
          variant="primary"
          icon={<TVLIcon className="w-6 h-6" />}
        />
        <StatsCard
          title="Total Pools"
          value={stats?.totalPools?.toLocaleString() || '-'}
          subtitle={`${stats?.totalProtocols || 0} protocols`}
          icon={<PoolsIcon className="w-6 h-6" />}
        />
        <StatsCard
          title="Average APY"
          value={`${stats?.averageApy?.toFixed(2) || 0}%`}
          subtitle={`Max: ${stats?.maxApy?.toFixed(2) || 0}%`}
          variant="success"
          icon={<APYIcon className="w-6 h-6" />}
        />
        <StatsCard
          title="Active Opportunities"
          value={stats?.activeOpportunities?.toString() || '-'}
          subtitle="Yield gaps & trending"
          variant="warning"
          icon={<OpportunitiesIcon className="w-6 h-6" />}
        />
      </div>

      {/* TVL by Chain */}
      {stats?.tvlByChain && (
        <div className="card">
          <h2 className="text-lg font-semibold text-white mb-4">TVL by Chain</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {Object.entries(stats.tvlByChain)
              .sort(([, a], [, b]) => b - a)
              .slice(0, 6)
              .map(([chain, tvl]) => (
                <div
                  key={chain}
                  className="p-4 bg-dark-800 rounded-lg cursor-pointer hover:bg-dark-700 transition-colors"
                  onClick={() => navigate(`/pools?chain=${chain}`)}
                >
                  <p className="text-sm text-gray-400 capitalize">{chain}</p>
                  <p className="text-lg font-semibold text-white mt-1">
                    {formatNumber(tvl)}
                  </p>
                  <p className="text-xs text-gray-500 mt-1">
                    {stats.poolsByChain?.[chain] || 0} pools
                  </p>
                </div>
              ))}
          </div>
        </div>
      )}

      {/* Top Pools */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-white">Top Pools</h2>
          <button
            className="btn-ghost text-sm text-primary-400"
            onClick={() => navigate('/pools')}
          >
            View All →
          </button>
        </div>
        {poolsLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="card h-48 animate-pulse bg-dark-800" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {poolsData?.data.map((pool) => (
              <PoolCard
                key={pool.id}
                pool={pool}
                onClick={() => navigate(`/pools/${encodeURIComponent(pool.id)}`)}
              />
            ))}
          </div>
        )}
      </div>

      {/* Active Opportunities */}
      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-white">
            Active Opportunities
          </h2>
          <button
            className="btn-ghost text-sm text-primary-400"
            onClick={() => navigate('/opportunities')}
          >
            View All →
          </button>
        </div>
        {oppsLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="card h-48 animate-pulse bg-dark-800" />
            ))}
          </div>
        ) : oppsData?.data.length === 0 ? (
          <div className="card text-center py-12">
            <p className="text-gray-400">No active opportunities found</p>
            <p className="text-sm text-gray-500 mt-2">
              Check back later for new yield opportunities
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {oppsData?.data.map((opportunity) => (
              <OpportunityCard
                key={opportunity.id}
                opportunity={opportunity}
                onClick={() => handleOpportunityClick(opportunity)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

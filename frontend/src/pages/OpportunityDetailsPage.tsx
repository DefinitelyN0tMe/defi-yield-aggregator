import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { opportunitiesApi, poolsApi } from '../services/api';
import { StatsCard, Sparkline, generateSparklineData, ExternalLinkIcon, RefreshIcon } from '../components';
import { getProtocolUrl, getExplorerUrl, defiLlamaLinks, getCoinGeckoUrl, getDexScreenerUrl } from '../utils/links';
import type { Opportunity } from '../types';

export function OpportunityDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  // Fetch opportunities list and find the one we want
  const { data: oppsData, isLoading, refetch: refetchOpps, isFetching: oppsFetching } = useQuery({
    queryKey: ['opportunities', 'all'],
    queryFn: () => opportunitiesApi.list({ limit: 100 }),
  });

  const opportunity = oppsData?.data.find((o) => o.id === id);

  // Fetch associated pool if available
  const poolId = opportunity?.poolId || opportunity?.targetPoolId || opportunity?.sourcePoolId;
  const { data: pool, refetch: refetchPool, isFetching: poolFetching } = useQuery({
    queryKey: ['pool', poolId],
    queryFn: () => poolsApi.get(poolId!),
    enabled: !!poolId,
  });

  // Combined fetching state for sync button
  const isSyncing = oppsFetching || poolFetching;

  // Sync all data
  const handleSync = () => {
    refetchOpps();
    if (poolId) refetchPool();
  };

  const formatNumber = (num: number, decimals = 2) => {
    if (num >= 1e9) return `$${(num / 1e9).toFixed(decimals)}B`;
    if (num >= 1e6) return `$${(num / 1e6).toFixed(decimals)}M`;
    if (num >= 1e3) return `$${(num / 1e3).toFixed(decimals)}K`;
    return `$${num.toFixed(decimals)}`;
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'yield-gap':
        return 'Yield Gap';
      case 'trending':
        return 'Trending';
      case 'high-score':
        return 'High Score';
      default:
        return type;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'yield-gap':
        return 'bg-purple-500/20 text-purple-400';
      case 'trending':
        return 'bg-blue-500/20 text-blue-400';
      case 'high-score':
        return 'bg-green-500/20 text-green-400';
      default:
        return 'bg-gray-500/20 text-gray-400';
    }
  };

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'low':
        return 'badge-success';
      case 'medium':
        return 'badge-warning';
      case 'high':
        return 'badge-danger';
      default:
        return 'badge-info';
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-400';
    if (score >= 60) return 'text-yellow-400';
    return 'text-red-400';
  };

  // Extract protocol from the opportunity
  const getProtocolFromOpportunity = (opp: Opportunity) => {
    const protocols = ['aave', 'compound', 'uniswap', 'curve', 'lido', 'gmx', 'pendle', 'velodrome', 'yearn', 'beefy', 'convex', 'rocket-pool', 'maker', 'sushiswap'];
    const descLower = opp.description.toLowerCase();
    const titleLower = opp.title.toLowerCase();
    for (const p of protocols) {
      if (descLower.includes(p) || titleLower.includes(p)) return p;
    }
    return pool?.protocol || null;
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="h-8 w-64 bg-dark-700 rounded animate-pulse" />
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-32 bg-dark-700 rounded-xl animate-pulse" />
          ))}
        </div>
        <div className="h-64 bg-dark-700 rounded-xl animate-pulse" />
      </div>
    );
  }

  if (!opportunity) {
    return (
      <div className="text-center py-12">
        <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-dark-800 flex items-center justify-center">
          <svg className="w-8 h-8 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 className="text-xl font-semibold text-white">Opportunity not found</h2>
        <p className="text-gray-400 mt-2">
          This opportunity may have expired or been removed.
        </p>
        <div className="flex items-center justify-center gap-4 mt-6">
          <button
            className="btn-primary"
            onClick={() => navigate('/opportunities')}
          >
            Browse Opportunities
          </button>
          <button
            className="btn-secondary"
            onClick={() => navigate(-1)}
          >
            Go Back
          </button>
        </div>
      </div>
    );
  }

  const protocol = getProtocolFromOpportunity(opportunity);
  const protocolUrl = protocol ? getProtocolUrl(protocol, opportunity.chain) : null;
  const explorerUrl = getExplorerUrl(opportunity.chain);

  return (
    <div className="space-y-6">
      {/* Back button and header */}
      <div className="flex items-start justify-between">
        <div>
          <button
            className="flex items-center gap-2 text-gray-400 hover:text-white mb-4 transition-colors"
            onClick={() => navigate('/opportunities')}
          >
            <span>←</span>
            <span>Back to Opportunities</span>
          </button>
          <div className="flex items-center gap-4 mb-2">
            <h1 className="text-3xl font-bold text-white">{opportunity.title}</h1>
          </div>
          <div className="flex items-center gap-2">
            <span className={`badge ${getTypeColor(opportunity.type)}`}>
              {getTypeLabel(opportunity.type)}
            </span>
            <span className={`badge ${getRiskColor(opportunity.riskLevel)}`}>
              {opportunity.riskLevel} risk
            </span>
            <span className="badge bg-dark-700 text-gray-300">
              {opportunity.chain}
            </span>
            {!opportunity.isActive && (
              <span className="badge badge-warning">Expired</span>
            )}
          </div>

          {/* External Links */}
          <div className="flex items-center gap-3 mt-4">
            {protocolUrl && (
              <a
                href={protocolUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-sm font-medium transition-colors"
              >
                <ExternalLinkIcon className="w-4 h-4" />
                Deposit on {protocol}
              </a>
            )}
            <a
              href={defiLlamaLinks.yields()}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              DefiLlama
            </a>
            <a
              href={explorerUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              Explorer
            </a>
            <a
              href={getDexScreenerUrl(opportunity.asset, opportunity.chain)}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              DexScreener
            </a>
            <a
              href={getCoinGeckoUrl(opportunity.asset)}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              CoinGecko
            </a>
          </div>
        </div>
        <div className="flex flex-col items-end gap-3">
          <button
            onClick={handleSync}
            disabled={isSyncing}
            className="btn-secondary flex items-center gap-2"
            title="Refresh opportunity data"
          >
            <RefreshIcon className={`w-4 h-4 ${isSyncing ? 'animate-spin' : ''}`} />
            {isSyncing ? 'Syncing...' : 'Sync'}
          </button>
          <div className={`text-4xl font-bold ${getScoreColor(opportunity.score)}`}>
            {opportunity.score.toFixed(0)}
            <span className="text-sm font-normal text-gray-400 ml-2">score</span>
          </div>
        </div>
      </div>

      {/* Description */}
      <div className="card">
        <p className="text-gray-300 text-lg">{opportunity.description}</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Current APY"
          value={`${opportunity.currentApy.toFixed(2)}%`}
          subtitle={`Asset: ${opportunity.asset}`}
          variant="success"
        />
        <StatsCard
          title={opportunity.type === 'yield-gap' ? 'APY Difference' : 'APY Growth'}
          value={`+${(opportunity.apyDifference || opportunity.apyGrowth || 0).toFixed(2)}%`}
          subtitle="Opportunity spread"
          variant="primary"
        />
        <StatsCard
          title="Total Value Locked"
          value={formatNumber(opportunity.tvl)}
          subtitle={`On ${opportunity.chain}`}
        />
        <StatsCard
          title="Potential Profit"
          value={`${opportunity.potentialProfit.toFixed(2)}%`}
          subtitle="Estimated gain"
          variant="warning"
        />
      </div>

      {/* APY Trend Chart */}
      {pool && (
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">APY Trend (7d)</h3>
          <div className="flex items-center gap-4">
            <Sparkline
              data={generateSparklineData(opportunity.currentApy, opportunity.apyGrowth || 5, 30)}
              width={400}
              height={80}
              showChange
            />
          </div>
        </div>
      )}

      {/* Associated Pool Info */}
      {pool && (
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">Associated Pool</h3>
          <div className="flex items-center justify-between p-4 bg-dark-800 rounded-lg">
            <div>
              <p className="text-white font-medium text-lg">{pool.symbol}</p>
              <div className="flex items-center gap-2 mt-1">
                <span className="badge bg-dark-700 text-gray-300">{pool.chain}</span>
                <span className="badge bg-dark-700 text-gray-300">{pool.protocol}</span>
              </div>
              <div className="flex items-center gap-4 mt-2 text-sm">
                <span className="text-gray-400">APY: <span className="text-green-400 font-medium">{pool.apy.toFixed(2)}%</span></span>
                <span className="text-gray-400">TVL: <span className="text-white font-medium">{formatNumber(pool.tvl)}</span></span>
                <span className="text-gray-400">Score: <span className={`font-medium ${getScoreColor(pool.score)}`}>{pool.score.toFixed(0)}</span></span>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <button
                className="btn-secondary"
                onClick={() => navigate(`/pools/${encodeURIComponent(pool.id)}`)}
              >
                View Pool Details
              </button>
              <a
                href={getProtocolUrl(pool.protocol, pool.chain)}
                target="_blank"
                rel="noopener noreferrer"
                className="btn-primary inline-flex items-center gap-2"
              >
                <span>Deposit</span>
                <ExternalLinkIcon className="w-4 h-4" />
              </a>
            </div>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="card bg-gradient-to-r from-primary-600/10 to-primary-500/5 border-primary-500/30">
        <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
        <div className="flex flex-wrap gap-3">
          {protocolUrl && (
            <a
              href={protocolUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="btn-primary inline-flex items-center gap-2"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              Deposit on {protocol}
            </a>
          )}
          {pool && (
            <button
              className="btn-secondary inline-flex items-center gap-2"
              onClick={() => navigate(`/pools/${encodeURIComponent(pool.id)}`)}
            >
              View Pool Details
            </button>
          )}
          <a
            href={defiLlamaLinks.yields()}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            Compare on DefiLlama
          </a>
          <a
            href={getDexScreenerUrl(opportunity.asset, opportunity.chain)}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            View on DexScreener
          </a>
          <a
            href={getCoinGeckoUrl(opportunity.asset)}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            View on CoinGecko
          </a>
        </div>
      </div>

      {/* Details */}
      <div className="card">
        <h3 className="text-lg font-semibold text-white mb-4">Details</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <p className="text-gray-400">Type</p>
            <p className="text-white mt-1">{getTypeLabel(opportunity.type)}</p>
          </div>
          <div>
            <p className="text-gray-400">Risk Level</p>
            <p className="text-white mt-1 capitalize">{opportunity.riskLevel}</p>
          </div>
          <div>
            <p className="text-gray-400">Detected</p>
            <p className="text-white mt-1">
              {new Date(opportunity.detectedAt).toLocaleString()}
            </p>
          </div>
          <div>
            <p className="text-gray-400">Status</p>
            <p className={`mt-1 ${opportunity.isActive ? 'text-green-400' : 'text-yellow-400'}`}>
              {opportunity.isActive ? 'Active' : 'Expired'}
            </p>
          </div>
        </div>
        {opportunity.expiresAt && (
          <div className="mt-4 pt-4 border-t border-dark-700">
            <p className="text-gray-400 text-sm">Expires At</p>
            <p className="text-white mt-1">
              {new Date(opportunity.expiresAt).toLocaleString()}
            </p>
          </div>
        )}
        <div className="mt-4 pt-4 border-t border-dark-700">
          <p className="text-gray-400 text-sm">Opportunity ID</p>
          <code className="text-xs text-gray-500 bg-dark-800 px-2 py-1 rounded mt-1 inline-block">
            {opportunity.id}
          </code>
        </div>
      </div>
    </div>
  );
}

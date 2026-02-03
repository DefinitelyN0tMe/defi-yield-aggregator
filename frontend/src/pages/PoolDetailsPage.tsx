import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { poolsApi } from '../services/api';
import { APYChart, StatsCard } from '../components';
import { getProtocolUrl, getExplorerUrl, defiLlamaLinks } from '../utils/links';

export function PoolDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [historyPeriod, setHistoryPeriod] = useState<'1h' | '24h' | '7d' | '30d'>('7d');

  // Decode the pool ID from URL
  const poolId = id ? decodeURIComponent(id) : '';

  // Fetch pool details
  const { data: pool, isLoading: poolLoading, error: poolError } = useQuery({
    queryKey: ['pool', poolId],
    queryFn: () => poolsApi.get(poolId),
    enabled: !!poolId,
    retry: false,
  });

  // Fetch pool history
  const { data: historyData, isLoading: historyLoading } = useQuery({
    queryKey: ['pool-history', poolId, historyPeriod],
    queryFn: () => poolsApi.getHistory(poolId, historyPeriod),
    enabled: !!poolId && !!pool,
  });

  const formatNumber = (num: number, decimals = 2) => {
    if (num >= 1e9) return `$${(num / 1e9).toFixed(decimals)}B`;
    if (num >= 1e6) return `$${(num / 1e6).toFixed(decimals)}M`;
    if (num >= 1e3) return `$${(num / 1e3).toFixed(decimals)}K`;
    return `$${num.toFixed(decimals)}`;
  };

  const getChangeColor = (change: number) => {
    if (change > 0) return 'text-green-400';
    if (change < 0) return 'text-red-400';
    return 'text-gray-400';
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-400';
    if (score >= 60) return 'text-yellow-400';
    return 'text-red-400';
  };

  if (poolLoading) {
    return (
      <div className="space-y-6">
        <div className="h-8 w-64 bg-dark-700 rounded animate-pulse" />
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-32 bg-dark-700 rounded-xl animate-pulse" />
          ))}
        </div>
        <div className="h-96 bg-dark-700 rounded-xl animate-pulse" />
      </div>
    );
  }

  if (poolError || !pool) {
    return (
      <div className="text-center py-12">
        <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-dark-800 flex items-center justify-center">
          <svg className="w-8 h-8 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 className="text-xl font-semibold text-white">Pool not found</h2>
        <p className="text-gray-400 mt-2">
          The pool you're looking for doesn't exist or has been removed.
        </p>
        <p className="text-gray-500 text-sm mt-1">
          Pool ID: {poolId}
        </p>
        <div className="flex items-center justify-center gap-4 mt-6">
          <button
            className="btn-primary"
            onClick={() => navigate('/pools')}
          >
            Browse All Pools
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

  const protocolUrl = getProtocolUrl(pool.protocol, pool.chain);
  const explorerUrl = getExplorerUrl(pool.chain);
  const defiLlamaUrl = defiLlamaLinks.protocol(pool.protocol);

  return (
    <div className="space-y-6">
      {/* Back button and header */}
      <div className="flex items-start justify-between">
        <div>
          <button
            className="flex items-center gap-2 text-gray-400 hover:text-white mb-4 transition-colors"
            onClick={() => navigate('/pools')}
          >
            <span>←</span>
            <span>Back to Pools</span>
          </button>
          <div className="flex items-center gap-4">
            <h1 className="text-3xl font-bold text-white">{pool.symbol}</h1>
            <div className="flex items-center gap-2">
              <span className="badge bg-dark-700 text-gray-300">
                {pool.chain}
              </span>
              <span className="badge bg-dark-700 text-gray-300">
                {pool.protocol}
              </span>
              {pool.stablecoin && (
                <span className="badge badge-info">Stablecoin</span>
              )}
            </div>
          </div>
          {pool.poolMeta && (
            <p className="text-gray-400 mt-2">{pool.poolMeta}</p>
          )}

          {/* External Links - Prominent */}
          <div className="flex flex-wrap items-center gap-3 mt-4">
            <a
              href={protocolUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white rounded-lg text-sm font-medium transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              Deposit on {pool.protocol}
            </a>
            <a
              href={defiLlamaUrl}
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
              href={`https://dexscreener.com/search?q=${pool.underlyingTokens[0] || pool.symbol}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              DexScreener
            </a>
            <a
              href={`https://www.coingecko.com/en/coins/${pool.underlyingTokens[0]?.toLowerCase() || pool.symbol.toLowerCase()}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-3 py-2 bg-dark-700 hover:bg-dark-600 text-gray-300 rounded-lg text-sm transition-colors"
            >
              <ExternalLinkIcon className="w-4 h-4" />
              CoinGecko
            </a>
          </div>
        </div>
        <div className={`text-4xl font-bold ${getScoreColor(pool.score)}`}>
          {pool.score.toFixed(0)}
          <span className="text-sm font-normal text-gray-400 ml-2">score</span>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Total APY"
          value={`${pool.apy.toFixed(2)}%`}
          subtitle={`Base: ${pool.apyBase.toFixed(2)}% | Reward: ${pool.apyReward.toFixed(2)}%`}
          variant="success"
        />
        <StatsCard
          title="Total Value Locked"
          value={formatNumber(pool.tvl)}
          subtitle={`Volume 24h: ${formatNumber(pool.volumeUsd1d)}`}
        />
        <StatsCard
          title="30d Average APY"
          value={`${pool.apyMean30d.toFixed(2)}%`}
          subtitle={`7d IL: ${pool.il7d.toFixed(2)}%`}
        />
        <StatsCard
          title="APY Changes"
          value={
            <div className="space-y-1 text-sm">
              <div className="flex items-center justify-between">
                <span className="text-gray-400">1h:</span>
                <span className={getChangeColor(pool.apyChange1h)}>
                  {pool.apyChange1h >= 0 ? '+' : ''}
                  {pool.apyChange1h.toFixed(2)}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-400">24h:</span>
                <span className={getChangeColor(pool.apyChange24h)}>
                  {pool.apyChange24h >= 0 ? '+' : ''}
                  {pool.apyChange24h.toFixed(2)}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-gray-400">7d:</span>
                <span className={getChangeColor(pool.apyChange7d)}>
                  {pool.apyChange7d >= 0 ? '+' : ''}
                  {pool.apyChange7d.toFixed(2)}%
                </span>
              </div>
            </div>
          }
        />
      </div>

      {/* APY Chart */}
      <APYChart
        data={historyData?.dataPoints || []}
        period={historyPeriod}
        onPeriodChange={setHistoryPeriod}
        loading={historyLoading}
        showTVL
      />

      {/* Token Info */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Underlying Tokens */}
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">
            Underlying Tokens
          </h3>
          <div className="space-y-2">
            {pool.underlyingTokens.length > 0 ? (
              pool.underlyingTokens.map((token, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between p-3 bg-dark-800 rounded-lg"
                >
                  <span className="text-white font-medium">{token}</span>
                  <a
                    href={`https://www.coingecko.com/en/coins/${token.toLowerCase()}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary-400 hover:text-primary-300 text-sm"
                  >
                    View on CoinGecko →
                  </a>
                </div>
              ))
            ) : (
              <p className="text-gray-400">No underlying tokens</p>
            )}
          </div>
        </div>

        {/* Reward Tokens */}
        <div className="card">
          <h3 className="text-lg font-semibold text-white mb-4">
            Reward Tokens
          </h3>
          <div className="space-y-2">
            {pool.rewardTokens.length > 0 ? (
              pool.rewardTokens.map((token, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between p-3 bg-dark-800 rounded-lg"
                >
                  <span className="text-primary-400 font-medium">{token}</span>
                  <a
                    href={`https://www.coingecko.com/en/coins/${token.toLowerCase()}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary-400 hover:text-primary-300 text-sm"
                  >
                    View on CoinGecko →
                  </a>
                </div>
              ))
            ) : (
              <p className="text-gray-400">No reward tokens (base yield only)</p>
            )}
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card bg-gradient-to-r from-primary-600/10 to-primary-500/5 border-primary-500/30">
        <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
        <div className="flex flex-wrap gap-3">
          <a
            href={protocolUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-primary inline-flex items-center gap-2"
          >
            <ExternalLinkIcon className="w-4 h-4" />
            Deposit on {pool.protocol}
          </a>
          <a
            href={`https://defillama.com/yields?chain=${pool.chain}&token=${pool.symbol}`}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            Compare on DefiLlama
          </a>
          <a
            href={`https://dexscreener.com/search?q=${pool.underlyingTokens[0] || pool.symbol}`}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            View on DexScreener
          </a>
          <a
            href={`https://www.coingecko.com/en/coins/${pool.underlyingTokens[0]?.toLowerCase() || pool.symbol.toLowerCase()}`}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            View on CoinGecko
          </a>
          <a
            href={explorerUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            {pool.chain} Explorer
          </a>
          <a
            href={`https://zapper.xyz/token/${pool.chain.toLowerCase()}/${pool.underlyingTokens[0]?.toLowerCase() || ''}`}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-secondary inline-flex items-center gap-2"
          >
            View on Zapper
          </a>
        </div>
      </div>

      {/* Metadata */}
      <div className="card">
        <h3 className="text-lg font-semibold text-white mb-4">Details</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <p className="text-gray-400">Exposure</p>
            <p className="text-white mt-1">{pool.exposure || '-'}</p>
          </div>
          <div>
            <p className="text-gray-400">7d Volume</p>
            <p className="text-white mt-1">{formatNumber(pool.volumeUsd7d)}</p>
          </div>
          <div>
            <p className="text-gray-400">Created</p>
            <p className="text-white mt-1">
              {new Date(pool.createdAt).toLocaleDateString()}
            </p>
          </div>
          <div>
            <p className="text-gray-400">Last Updated</p>
            <p className="text-white mt-1">
              {new Date(pool.updatedAt).toLocaleString()}
            </p>
          </div>
        </div>
        <div className="mt-4 pt-4 border-t border-dark-700">
          <p className="text-gray-400 text-sm">Pool ID</p>
          <code className="text-xs text-gray-500 bg-dark-800 px-2 py-1 rounded mt-1 inline-block">
            {pool.id}
          </code>
        </div>
      </div>
    </div>
  );
}

function ExternalLinkIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
    </svg>
  );
}

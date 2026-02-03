import type { Pool } from '../types';
import { getProtocolUrl, getExplorerUrl } from '../utils/links';
import { formatNumber, getScoreColor, getChangeColor, safeFirstLower } from '../utils/format';
import { Sparkline, generateSparklineData } from './Sparkline';
import { ExternalLinkIcon } from './Icons';

interface PoolCardProps {
  pool: Pool;
  onClick?: () => void;
}

export function PoolCard({ pool, onClick }: PoolCardProps) {
  const protocolUrl = getProtocolUrl(pool.protocol, pool.chain);
  const explorerUrl = getExplorerUrl(pool.chain);

  const handleExternalClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  return (
    <div
      className="card-hover cursor-pointer"
      onClick={onClick}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div>
          <h3 className="text-lg font-semibold text-white">{pool.symbol}</h3>
          <div className="flex items-center gap-2 mt-1">
            <span className="badge bg-dark-700 text-gray-300">{pool.chain}</span>
            <span className="badge bg-dark-700 text-gray-300">{pool.protocol}</span>
          </div>
        </div>
        <div className={`text-2xl font-bold ${getScoreColor(pool.score)}`}>
          {pool.score.toFixed(0)}
        </div>
      </div>

      {/* APY and Chart */}
      <div className="flex items-center justify-between mb-4">
        <div>
          <p className="text-sm text-gray-400">APY</p>
          <p className="text-2xl font-bold text-green-400">
            {pool.apy.toFixed(2)}%
          </p>
          <p className={`text-xs ${getChangeColor(pool.apyChange24h)}`}>
            {pool.apyChange24h >= 0 ? '↑' : '↓'} {Math.abs(pool.apyChange24h).toFixed(2)}% 24h
          </p>
        </div>
        <div className="text-right">
          <p className="text-xs text-gray-500 mb-1">7d Trend</p>
          <Sparkline
            data={generateSparklineData(pool.apy, pool.apyChange7d, 14)}
            width={80}
            height={32}
          />
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <p className="text-sm text-gray-400">TVL</p>
          <p className="text-lg font-semibold text-white">
            {formatNumber(pool.tvl)}
          </p>
        </div>
        <div>
          <p className="text-sm text-gray-400">Volume 24h</p>
          <p className="text-lg font-semibold text-white">
            {formatNumber(pool.volumeUsd1d)}
          </p>
        </div>
      </div>

      {/* APY Breakdown */}
      <div className="pt-3 border-t border-dark-700 mb-3">
        <div className="flex items-center justify-between text-sm">
          <div className="flex items-center gap-4">
            <div>
              <span className="text-gray-400">Base: </span>
              <span className="text-white">{pool.apyBase.toFixed(2)}%</span>
            </div>
            <div>
              <span className="text-gray-400">Reward: </span>
              <span className="text-primary-400">{pool.apyReward.toFixed(2)}%</span>
            </div>
          </div>
          {pool.stablecoin && (
            <span className="badge badge-info">Stablecoin</span>
          )}
        </div>
      </div>

      {/* Quick Links */}
      <div className="pt-3 border-t border-dark-700" onClick={handleExternalClick}>
        <div className="flex items-center gap-2">
          <a
            href={protocolUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="flex-1 inline-flex items-center justify-center gap-1 px-3 py-2 text-sm font-medium text-white bg-primary-600 hover:bg-primary-500 rounded-lg transition-colors"
          >
            <span>Deposit on {pool.protocol}</span>
            <ExternalLinkIcon className="w-3 h-3" />
          </a>
        </div>
        <div className="flex items-center justify-center gap-4 mt-2">
          <a
            href={`https://defillama.com/yields?project=${pool.protocol}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-gray-400 hover:text-primary-400"
          >
            DefiLlama
          </a>
          <span className="text-gray-600">•</span>
          <a
            href={explorerUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-gray-400 hover:text-primary-400"
          >
            {pool.chain} Explorer
          </a>
          <span className="text-gray-600">•</span>
          <a
            href={`https://www.coingecko.com/en/coins/${safeFirstLower(pool.underlyingTokens, pool.symbol)}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-gray-400 hover:text-primary-400"
          >
            CoinGecko
          </a>
        </div>
      </div>
    </div>
  );
}

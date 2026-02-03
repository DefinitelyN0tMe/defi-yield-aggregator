import type { Opportunity } from '../types';
import { getProtocolUrl, getExplorerUrl } from '../utils/links';
import { formatNumber, getTypeLabel, getTypeColor, getRiskColor, formatTimeAgo } from '../utils/format';
import { ExternalLinkIcon } from './Icons';
import { KNOWN_PROTOCOLS } from '../utils/constants';

interface OpportunityCardProps {
  opportunity: Opportunity;
  onClick?: () => void;
}

export function OpportunityCard({ opportunity, onClick }: OpportunityCardProps) {
  // Extract protocol from the opportunity (parse from description or use a default)
  const getProtocolFromOpportunity = () => {
    const descLower = opportunity.description.toLowerCase();
    for (const p of KNOWN_PROTOCOLS) {
      if (descLower.includes(p)) return p;
    }
    return null;
  };

  const protocol = getProtocolFromOpportunity();
  const explorerUrl = getExplorerUrl(opportunity.chain);

  const handleExternalClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  return (
    <div
      className="card-hover cursor-pointer"
      onClick={onClick}
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className={`badge ${getTypeColor(opportunity.type)}`}>
            {getTypeLabel(opportunity.type)}
          </span>
          <span className={`badge ${getRiskColor(opportunity.riskLevel)}`}>
            {opportunity.riskLevel} risk
          </span>
        </div>
        <span className="text-xs text-gray-500">
          {formatTimeAgo(opportunity.detectedAt)}
        </span>
      </div>

      <h3 className="text-lg font-semibold text-white mb-2">
        {opportunity.title}
      </h3>
      <p className="text-sm text-gray-400 mb-4 line-clamp-2">
        {opportunity.description}
      </p>

      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <p className="text-xs text-gray-500">Current APY</p>
          <p className="text-lg font-semibold text-green-400">
            {opportunity.currentApy.toFixed(2)}%
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-500">
            {opportunity.type === 'yield-gap' ? 'APY Diff' : 'APY Growth'}
          </p>
          <p className="text-lg font-semibold text-primary-400">
            +{(opportunity.apyDifference || opportunity.apyGrowth).toFixed(2)}%
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-500">TVL</p>
          <p className="text-lg font-semibold text-white">
            {formatNumber(opportunity.tvl)}
          </p>
        </div>
      </div>

      <div className="flex items-center justify-between pt-4 border-t border-dark-700">
        <div className="flex items-center gap-2">
          <span className="badge bg-dark-700 text-gray-300">
            {opportunity.chain}
          </span>
          <span className="text-sm text-gray-400">{opportunity.asset}</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-400">Score:</span>
          <span className={`font-bold ${
            opportunity.score >= 80
              ? 'text-green-400'
              : opportunity.score >= 60
              ? 'text-yellow-400'
              : 'text-red-400'
          }`}>
            {opportunity.score.toFixed(0)}
          </span>
        </div>
      </div>

      {/* Quick Links */}
      <div className="pt-3 border-t border-dark-700 mt-3" onClick={handleExternalClick}>
        <div className="flex items-center gap-2">
          {protocol && (
            <a
              href={getProtocolUrl(protocol)}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 inline-flex items-center justify-center gap-1 px-3 py-2 text-sm font-medium text-white bg-primary-600 hover:bg-primary-500 rounded-lg transition-colors"
            >
              <span>Deposit on {protocol}</span>
              <ExternalLinkIcon className="w-3 h-3" />
            </a>
          )}
          {!protocol && (
            <a
              href={`https://defillama.com/yields?chain=${opportunity.chain}`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex-1 inline-flex items-center justify-center gap-1 px-3 py-2 text-sm font-medium text-white bg-primary-600 hover:bg-primary-500 rounded-lg transition-colors"
            >
              <span>Find on DefiLlama</span>
              <ExternalLinkIcon className="w-3 h-3" />
            </a>
          )}
        </div>
        <div className="flex items-center justify-center gap-4 mt-2">
          <a
            href={`https://defillama.com/yields?chain=${opportunity.chain}`}
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
            {opportunity.chain} Explorer
          </a>
          <span className="text-gray-600">•</span>
          <a
            href={`https://dexscreener.com/search?q=${opportunity.asset}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs text-gray-400 hover:text-primary-400"
          >
            DexScreener
          </a>
        </div>
      </div>

      {!opportunity.isActive && (
        <div className="mt-3 pt-3 border-t border-dark-700">
          <span className="badge badge-warning">Expired</span>
        </div>
      )}
    </div>
  );
}

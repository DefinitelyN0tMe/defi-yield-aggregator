import { useState } from 'react';
import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { opportunitiesApi } from '../services/api';
import { OpportunityCard } from '../components';
import { useOpportunityAlerts } from '../hooks/useWebSocket';
import { getProtocolUrl } from '../utils/links';
import type { OpportunityFilter, OpportunityType, RiskLevel, Opportunity } from '../types';

export function OpportunitiesPage() {
  const navigate = useNavigate();
  const [filter, setFilter] = useState<OpportunityFilter>({
    activeOnly: true,
    sortBy: 'score',
    sortOrder: 'desc',
    limit: 20,
    offset: 0,
  });

  // Fetch opportunities
  const { data: oppsData, isLoading, refetch } = useQuery({
    queryKey: ['opportunities', filter],
    queryFn: () => opportunitiesApi.list(filter),
    placeholderData: keepPreviousData,
  });

  // Fetch trending pools
  const { data: trendingData } = useQuery({
    queryKey: ['trending'],
    queryFn: () => opportunitiesApi.getTrending({ limit: 5 }),
  });

  // Subscribe to real-time opportunity alerts
  useOpportunityAlerts(() => {
    // New opportunity received - refresh the list
    refetch();
  });

  const handleFilterChange = (newFilter: Partial<OpportunityFilter>) => {
    setFilter((prev) => ({ ...prev, ...newFilter, offset: 0 }));
  };

  const handleLoadMore = () => {
    setFilter((prev) => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 20),
    }));
  };

  // Navigate to opportunity details page
  const handleOpportunityClick = (opportunity: Opportunity) => {
    navigate(`/opportunities/${encodeURIComponent(opportunity.id)}`);
  };

  const typeOptions: { value: OpportunityType | ''; label: string }[] = [
    { value: '', label: 'All Types' },
    { value: 'yield-gap', label: 'Yield Gap' },
    { value: 'trending', label: 'Trending' },
    { value: 'high-score', label: 'High Score' },
  ];

  const riskOptions: { value: RiskLevel | ''; label: string }[] = [
    { value: '', label: 'All Risk Levels' },
    { value: 'low', label: 'Low Risk' },
    { value: 'medium', label: 'Medium Risk' },
    { value: 'high', label: 'High Risk' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-white">Opportunities</h1>
        <p className="mt-2 text-gray-400">
          Discover yield farming opportunities based on APY gaps, trending pools, and
          risk-adjusted scores
        </p>
      </div>

      {/* Trending Pools */}
      {trendingData?.data && trendingData.data.length > 0 && (
        <div className="card bg-gradient-to-r from-primary-600/10 to-primary-500/5 border-primary-500/30">
          <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            <TrendingIcon className="w-5 h-5 text-primary-400" />
            Trending Now
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            {trendingData.data.map((trending, index) => (
              <div
                key={trending.pool.id}
                className="p-4 bg-dark-800/50 rounded-lg cursor-pointer hover:bg-dark-700/50 transition-colors"
                onClick={() => navigate(`/pools/${encodeURIComponent(trending.pool.id)}`)}
              >
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-2xl font-bold text-primary-400">
                    #{index + 1}
                  </span>
                </div>
                <p className="font-medium text-white truncate">
                  {trending.pool.symbol}
                </p>
                <p className="text-sm text-gray-400">{trending.pool.protocol}</p>
                <div className="mt-2">
                  <span className="text-green-400 font-semibold">
                    {trending.pool.apy.toFixed(2)}%
                  </span>
                  <span className="text-sm text-gray-500 ml-2">APY</span>
                </div>
                <div className="text-sm text-primary-400 mt-1">
                  +{trending.apyGrowth24h.toFixed(2)}% 24h
                </div>
                <a
                  href={getProtocolUrl(trending.pool.protocol, trending.pool.chain)}
                  target="_blank"
                  rel="noopener noreferrer"
                  onClick={(e) => e.stopPropagation()}
                  className="inline-flex items-center gap-1 text-xs text-gray-400 hover:text-primary-400 mt-2"
                >
                  Open protocol →
                </a>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-4">
        <select
          className="select min-w-[150px]"
          value={filter.type || ''}
          onChange={(e) =>
            handleFilterChange({
              type: (e.target.value as OpportunityType) || undefined,
            })
          }
        >
          {typeOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>

        <select
          className="select min-w-[150px]"
          value={filter.riskLevel || ''}
          onChange={(e) =>
            handleFilterChange({
              riskLevel: (e.target.value as RiskLevel) || undefined,
            })
          }
        >
          {riskOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>

        <select
          className="select min-w-[150px]"
          value={filter.sortBy || 'score'}
          onChange={(e) =>
            handleFilterChange({
              sortBy: e.target.value as 'score' | 'profit' | 'apy' | 'detectedAt',
            })
          }
        >
          <option value="score">Sort by Score</option>
          <option value="profit">Sort by Profit</option>
          <option value="apy">Sort by APY</option>
          <option value="detectedAt">Sort by Date</option>
        </select>

        <label className="flex items-center gap-2 text-sm text-gray-400">
          <input
            type="checkbox"
            checked={filter.activeOnly || false}
            onChange={(e) =>
              handleFilterChange({ activeOnly: e.target.checked })
            }
            className="rounded bg-dark-700 border-dark-600 text-primary-500 focus:ring-primary-500"
          />
          Active only
        </label>

        <div className="flex-1" />

        <div className="text-sm text-gray-400">
          {oppsData?.total || 0} opportunities found
        </div>
      </div>

      {/* Opportunities Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="card h-48 animate-pulse bg-dark-800" />
          ))}
        </div>
      ) : oppsData?.data.length === 0 ? (
        <div className="card text-center py-12">
          <NoDataIcon className="w-16 h-16 text-gray-600 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-white">
            No opportunities found
          </h3>
          <p className="text-gray-400 mt-2">
            Try adjusting your filters or check back later
          </p>
          <button
            className="btn-primary mt-4"
            onClick={() =>
              setFilter({
                activeOnly: true,
                sortBy: 'score',
                sortOrder: 'desc',
                limit: 20,
                offset: 0,
              })
            }
          >
            Reset Filters
          </button>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {oppsData?.data.map((opportunity) => (
              <OpportunityCard
                key={opportunity.id}
                opportunity={opportunity}
                onClick={() => handleOpportunityClick(opportunity)}
              />
            ))}
          </div>

          {/* Load More */}
          {oppsData?.hasMore && (
            <div className="text-center">
              <button className="btn-secondary" onClick={handleLoadMore}>
                Load More
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}

function TrendingIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
      />
    </svg>
  );
}

function NoDataIcon({ className }: { className?: string }) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      />
    </svg>
  );
}

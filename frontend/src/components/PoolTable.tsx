import { useState, useEffect, useRef } from 'react';
import type { Pool, PoolFilter } from '../types';
import { getProtocolUrl, getExplorerUrl, defiLlamaLinks } from '../utils/links';
import { formatNumber, getScoreColor, getChangeColor } from '../utils/format';
import { Sparkline, generateSparklineData } from './Sparkline';
import { ExternalLinkIcon } from './Icons';

interface PoolTableProps {
  pools: Pool[];
  loading?: boolean;
  onPoolClick?: (pool: Pool) => void;
  onFilterChange?: (filter: PoolFilter) => void;
  filter?: PoolFilter;
  total?: number;
  hasMore?: boolean;
  onLoadMore?: () => void;
}

export function PoolTable({
  pools,
  loading,
  onPoolClick,
  onFilterChange,
  filter = {},
  total,
  hasMore,
  onLoadMore,
}: PoolTableProps) {
  // Local search input state with debounce
  const [searchInput, setSearchInput] = useState(filter.search || '');
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  // Update local state when filter.search changes externally
  useEffect(() => {
    setSearchInput(filter.search || '');
  }, [filter.search]);

  // Debounced search - waits 300ms after user stops typing
  const handleSearchChange = (value: string) => {
    setSearchInput(value);

    // Clear previous timeout
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    // Set new timeout to update filter after 300ms
    debounceRef.current = setTimeout(() => {
      onFilterChange?.({ ...filter, search: value || undefined, offset: 0 });
    }, 300);
  };

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, []);

  const handleSort = (sortBy: 'apy' | 'tvl' | 'score') => {
    const newOrder =
      filter.sortBy === sortBy && filter.sortOrder === 'desc' ? 'asc' : 'desc';
    onFilterChange?.({ ...filter, sortBy, sortOrder: newOrder });
  };

  const SortIcon = ({ field }: { field: string }) => {
    if (filter.sortBy !== field) return null;
    return (
      <span className="ml-1">
        {filter.sortOrder === 'asc' ? '↑' : '↓'}
      </span>
    );
  };

  const handleExternalClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  return (
    <div className="card p-0 overflow-hidden">
      {/* Filters */}
      <div className="p-4 border-b border-dark-700">
        <div className="flex flex-wrap items-center gap-4">
          <div className="flex-1 min-w-[200px]">
            <input
              type="text"
              placeholder="Search all pools by symbol, protocol, or chain..."
              className="input"
              value={searchInput}
              onChange={(e) => handleSearchChange(e.target.value)}
            />
          </div>
          <div className="flex items-center gap-2">
            <select
              className="select"
              value={filter.sortBy || 'score'}
              onChange={(e) =>
                onFilterChange?.({
                  ...filter,
                  sortBy: e.target.value as 'apy' | 'tvl' | 'score',
                })
              }
            >
              <option value="score">Sort by Score</option>
              <option value="apy">Sort by APY</option>
              <option value="tvl">Sort by TVL</option>
            </select>
            <label className="flex items-center gap-2 text-sm text-gray-400">
              <input
                type="checkbox"
                checked={filter.stablecoin || false}
                onChange={(e) =>
                  onFilterChange?.({
                    ...filter,
                    stablecoin: e.target.checked || undefined,
                  })
                }
                className="rounded bg-dark-700 border-dark-600 text-primary-500 focus:ring-primary-500"
              />
              Stablecoins only
            </label>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-dark-800">
            <tr>
              <th className="table-header">Pool</th>
              <th className="table-header">Chain / Protocol</th>
              <th
                className="table-header cursor-pointer hover:text-white transition-colors"
                onClick={() => handleSort('apy')}
              >
                APY <SortIcon field="apy" />
              </th>
              <th className="table-header">7d Trend</th>
              <th
                className="table-header cursor-pointer hover:text-white transition-colors"
                onClick={() => handleSort('tvl')}
              >
                TVL <SortIcon field="tvl" />
              </th>
              <th
                className="table-header cursor-pointer hover:text-white transition-colors"
                onClick={() => handleSort('score')}
              >
                Score <SortIcon field="score" />
              </th>
              <th className="table-header">Quick Links</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-dark-700">
            {loading ? (
              Array.from({ length: 10 }).map((_, i) => (
                <tr key={i}>
                  <td className="table-cell" colSpan={7}>
                    <div className="h-12 bg-dark-700 rounded animate-pulse" />
                  </td>
                </tr>
              ))
            ) : pools.length === 0 ? (
              <tr>
                <td
                  className="table-cell text-center text-gray-400 py-12"
                  colSpan={7}
                >
                  <div className="flex flex-col items-center">
                    <svg className="w-12 h-12 text-gray-600 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <p className="text-lg font-medium">No pools found</p>
                    <p className="text-sm text-gray-500 mt-1">Try adjusting your search or filters</p>
                  </div>
                </td>
              </tr>
            ) : (
              pools.map((pool) => (
                <tr
                  key={pool.id}
                  className="hover:bg-dark-800/50 cursor-pointer transition-colors"
                  onClick={() => onPoolClick?.(pool)}
                >
                  <td className="table-cell">
                    <div className="flex items-center gap-2">
                      <div>
                        <span className="font-medium text-white block">
                          {pool.symbol}
                        </span>
                        <span className="text-xs text-gray-500">
                          {pool.poolMeta?.slice(0, 30) || `${pool.protocol} Pool`}
                        </span>
                      </div>
                      {pool.stablecoin && (
                        <span className="badge badge-info text-xs">Stable</span>
                      )}
                    </div>
                  </td>
                  <td className="table-cell">
                    <div className="flex flex-col gap-1">
                      <span className="badge bg-dark-700 text-gray-300 w-fit">
                        {pool.chain}
                      </span>
                      <span className="text-sm text-gray-400">{pool.protocol}</span>
                    </div>
                  </td>
                  <td className="table-cell">
                    <div>
                      <span className="text-green-400 font-bold text-lg">
                        {pool.apy.toFixed(2)}%
                      </span>
                      <div className="text-xs text-gray-500 mt-1">
                        <span className={getChangeColor(pool.apyChange24h)}>
                          {pool.apyChange24h >= 0 ? '↑' : '↓'} {Math.abs(pool.apyChange24h).toFixed(2)}%
                        </span>
                        <span className="text-gray-600 mx-1">24h</span>
                      </div>
                    </div>
                  </td>
                  <td className="table-cell">
                    <Sparkline
                      data={generateSparklineData(pool.apy, pool.apyChange7d, 14)}
                      width={70}
                      height={28}
                    />
                  </td>
                  <td className="table-cell">
                    <div>
                      <span className="font-medium text-white">
                        {formatNumber(pool.tvl)}
                      </span>
                      <div className="text-xs text-gray-500">
                        Vol: {formatNumber(pool.volumeUsd1d)}
                      </div>
                    </div>
                  </td>
                  <td className="table-cell">
                    <div className={`text-xl font-bold ${getScoreColor(pool.score)}`}>
                      {pool.score.toFixed(0)}
                    </div>
                  </td>
                  <td className="table-cell">
                    <div className="flex flex-col gap-1" onClick={handleExternalClick}>
                      <a
                        href={getProtocolUrl(pool.protocol, pool.chain)}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium text-white bg-primary-600 hover:bg-primary-500 rounded transition-colors"
                      >
                        <span>Deposit</span>
                        <ExternalLinkIcon className="w-3 h-3" />
                      </a>
                      <div className="flex items-center gap-2">
                        <a
                          href={defiLlamaLinks.protocol(pool.protocol)}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-xs text-gray-400 hover:text-primary-400"
                        >
                          DefiLlama
                        </a>
                        <span className="text-gray-600">•</span>
                        <a
                          href={getExplorerUrl(pool.chain)}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-xs text-gray-400 hover:text-primary-400"
                        >
                          Explorer
                        </a>
                      </div>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Footer */}
      <div className="p-4 border-t border-dark-700 flex items-center justify-between">
        <p className="text-sm text-gray-400">
          Showing {pools.length} of {total || pools.length} pools
        </p>
        <div className="flex items-center gap-3">
          {hasMore && (
            <button className="btn-secondary" onClick={onLoadMore}>
              Load More
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

import { useState } from 'react';
import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { poolsApi, chainsApi } from '../services/api';
import { PoolTable } from '../components';
import type { PoolFilter } from '../types';

export function PoolsPage() {
  const navigate = useNavigate();
  const [filter, setFilter] = useState<PoolFilter>({
    sortBy: 'score',
    sortOrder: 'desc',
    limit: 25,
    offset: 0,
  });

  // Fetch pools
  const { data: poolsData, isLoading } = useQuery({
    queryKey: ['pools', filter],
    queryFn: () => poolsApi.list(filter),
    placeholderData: keepPreviousData,
  });

  // Fetch chains for filter
  const { data: chainsData } = useQuery({
    queryKey: ['chains'],
    queryFn: chainsApi.list,
  });

  const handleFilterChange = (newFilter: PoolFilter) => {
    setFilter({ ...newFilter, offset: 0 });
  };

  const handleLoadMore = () => {
    setFilter((prev) => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 25),
    }));
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-white">Pools</h1>
          <p className="mt-2 text-gray-400">
            Browse and filter yield farming pools across all chains
          </p>
        </div>
        <div className="flex items-center gap-4">
          {/* Chain filter */}
          <select
            className="select min-w-[150px]"
            value={filter.chain || ''}
            onChange={(e) =>
              handleFilterChange({
                ...filter,
                chain: e.target.value || undefined,
              })
            }
          >
            <option value="">All Chains</option>
            {chainsData?.data.map((chain) => (
              <option key={chain.name} value={chain.name}>
                {chain.displayName} ({chain.poolCount})
              </option>
            ))}
          </select>

          {/* APY Range */}
          <div className="hidden md:flex items-center gap-2">
            <input
              type="number"
              placeholder="Min APY"
              className="input w-24"
              value={filter.minApy || ''}
              onChange={(e) =>
                handleFilterChange({
                  ...filter,
                  minApy: e.target.value ? Number(e.target.value) : undefined,
                })
              }
            />
            <span className="text-gray-500">-</span>
            <input
              type="number"
              placeholder="Max APY"
              className="input w-24"
              value={filter.maxApy || ''}
              onChange={(e) =>
                handleFilterChange({
                  ...filter,
                  maxApy: e.target.value ? Number(e.target.value) : undefined,
                })
              }
            />
          </div>

          {/* Min TVL */}
          <select
            className="select min-w-[120px]"
            value={filter.minTvl || ''}
            onChange={(e) =>
              handleFilterChange({
                ...filter,
                minTvl: e.target.value ? Number(e.target.value) : undefined,
              })
            }
          >
            <option value="">Any TVL</option>
            <option value="100000">$100K+</option>
            <option value="1000000">$1M+</option>
            <option value="10000000">$10M+</option>
            <option value="100000000">$100M+</option>
          </select>
        </div>
      </div>

      {/* Active filters */}
      {(filter.chain || filter.minApy || filter.maxApy || filter.minTvl || filter.stablecoin) && (
        <div className="flex items-center gap-2 flex-wrap">
          <span className="text-sm text-gray-400">Active filters:</span>
          {filter.chain && (
            <FilterBadge
              label={`Chain: ${filter.chain}`}
              onRemove={() => handleFilterChange({ ...filter, chain: undefined })}
            />
          )}
          {filter.minApy && (
            <FilterBadge
              label={`Min APY: ${filter.minApy}%`}
              onRemove={() => handleFilterChange({ ...filter, minApy: undefined })}
            />
          )}
          {filter.maxApy && (
            <FilterBadge
              label={`Max APY: ${filter.maxApy}%`}
              onRemove={() => handleFilterChange({ ...filter, maxApy: undefined })}
            />
          )}
          {filter.minTvl && (
            <FilterBadge
              label={`Min TVL: $${(filter.minTvl / 1e6).toFixed(0)}M`}
              onRemove={() => handleFilterChange({ ...filter, minTvl: undefined })}
            />
          )}
          {filter.stablecoin && (
            <FilterBadge
              label="Stablecoins only"
              onRemove={() => handleFilterChange({ ...filter, stablecoin: undefined })}
            />
          )}
          <button
            className="text-sm text-primary-400 hover:text-primary-300"
            onClick={() =>
              setFilter({
                sortBy: 'score',
                sortOrder: 'desc',
                limit: 25,
                offset: 0,
              })
            }
          >
            Clear all
          </button>
        </div>
      )}

      {/* Pool Table */}
      <PoolTable
        pools={poolsData?.data || []}
        loading={isLoading}
        onPoolClick={(pool) => navigate(`/pools/${pool.id}`)}
        onFilterChange={handleFilterChange}
        filter={filter}
        total={poolsData?.total}
        hasMore={poolsData?.hasMore}
        onLoadMore={handleLoadMore}
      />
    </div>
  );
}

function FilterBadge({
  label,
  onRemove,
}: {
  label: string;
  onRemove: () => void;
}) {
  return (
    <span className="inline-flex items-center gap-1 px-2 py-1 bg-dark-700 rounded-lg text-sm text-gray-300">
      {label}
      <button
        className="text-gray-400 hover:text-white"
        onClick={onRemove}
      >
        ×
      </button>
    </span>
  );
}

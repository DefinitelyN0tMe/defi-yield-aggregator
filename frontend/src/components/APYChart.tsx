import {
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  AreaChart,
} from 'recharts';
import type { HistoricalAPY } from '../types';

interface APYChartProps {
  data: HistoricalAPY[];
  period: '1h' | '24h' | '7d' | '30d';
  onPeriodChange?: (period: '1h' | '24h' | '7d' | '30d') => void;
  loading?: boolean;
  showTVL?: boolean;
}

export function APYChart({
  data,
  period,
  onPeriodChange,
  loading,
  showTVL = false,
}: APYChartProps) {
  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    if (period === '1h') {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    if (period === '24h') {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
  };

  const formatNumber = (num: number) => {
    if (num >= 1e9) return `$${(num / 1e9).toFixed(1)}B`;
    if (num >= 1e6) return `$${(num / 1e6).toFixed(1)}M`;
    if (num >= 1e3) return `$${(num / 1e3).toFixed(1)}K`;
    return `$${num.toFixed(0)}`;
  };

  const periods: { value: '1h' | '24h' | '7d' | '30d'; label: string }[] = [
    { value: '1h', label: '1H' },
    { value: '24h', label: '24H' },
    { value: '7d', label: '7D' },
    { value: '30d', label: '30D' },
  ];

  const chartData = data.map((d) => ({
    ...d,
    time: formatDate(d.timestamp),
  }));

  if (loading) {
    return (
      <div className="card">
        <div className="h-64 flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-2 border-primary-500 border-t-transparent" />
        </div>
      </div>
    );
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-white">APY History</h3>
        <div className="flex items-center gap-1">
          {periods.map((p) => (
            <button
              key={p.value}
              className={`px-3 py-1 text-sm rounded-lg transition-colors ${
                period === p.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-dark-700 text-gray-400 hover:text-white'
              }`}
              onClick={() => onPeriodChange?.(p.value)}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      {chartData.length === 0 ? (
        <div className="h-64 flex items-center justify-center text-gray-400">
          No data available for this period
        </div>
      ) : (
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="apyGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#22c55e" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#22c55e" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="tvlGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#6366f1" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid
                strokeDasharray="3 3"
                stroke="#374151"
                vertical={false}
              />
              <XAxis
                dataKey="time"
                stroke="#6b7280"
                tick={{ fill: '#9ca3af', fontSize: 12 }}
                tickLine={false}
                axisLine={false}
              />
              <YAxis
                yAxisId="apy"
                stroke="#6b7280"
                tick={{ fill: '#9ca3af', fontSize: 12 }}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => `${value.toFixed(1)}%`}
              />
              {showTVL && (
                <YAxis
                  yAxisId="tvl"
                  orientation="right"
                  stroke="#6b7280"
                  tick={{ fill: '#9ca3af', fontSize: 12 }}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={formatNumber}
                />
              )}
              <Tooltip
                contentStyle={{
                  backgroundColor: '#1f2937',
                  border: '1px solid #374151',
                  borderRadius: '0.5rem',
                }}
                labelStyle={{ color: '#9ca3af' }}
                itemStyle={{ color: '#fff' }}
                formatter={(value: number, name: string) => {
                  if (name === 'apy') return [`${value.toFixed(2)}%`, 'APY'];
                  if (name === 'tvl') return [formatNumber(value), 'TVL'];
                  return [value, name];
                }}
              />
              <Area
                yAxisId="apy"
                type="monotone"
                dataKey="apy"
                stroke="#22c55e"
                strokeWidth={2}
                fill="url(#apyGradient)"
              />
              {showTVL && (
                <Area
                  yAxisId="tvl"
                  type="monotone"
                  dataKey="tvl"
                  stroke="#6366f1"
                  strokeWidth={2}
                  fill="url(#tvlGradient)"
                />
              )}
            </AreaChart>
          </ResponsiveContainer>
        </div>
      )}

      {chartData.length > 0 && (
        <div className="mt-4 pt-4 border-t border-dark-700 flex items-center gap-6 text-sm">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full bg-green-500" />
            <span className="text-gray-400">Total APY</span>
          </div>
          {showTVL && (
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded-full bg-primary-500" />
              <span className="text-gray-400">TVL</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

import { useMemo } from 'react';

interface SparklineProps {
  data: number[];
  width?: number;
  height?: number;
  color?: string;
  showChange?: boolean;
}

export function Sparkline({
  data,
  width = 80,
  height = 24,
  color,
  showChange = false,
}: SparklineProps) {
  const { path, changePercent, isPositive } = useMemo(() => {
    if (!data || data.length < 2) {
      return { path: '', changePercent: 0, isPositive: true };
    }

    const min = Math.min(...data);
    const max = Math.max(...data);
    const range = max - min || 1;

    const points = data.map((value, index) => {
      const x = (index / (data.length - 1)) * width;
      const y = height - ((value - min) / range) * height;
      return `${x},${y}`;
    });

    const first = data[0];
    const last = data[data.length - 1];
    const change = ((last - first) / first) * 100;

    return {
      path: `M${points.join(' L')}`,
      changePercent: change,
      isPositive: change >= 0,
    };
  }, [data, width, height]);

  const strokeColor = color || (isPositive ? '#22c55e' : '#ef4444');

  if (!data || data.length < 2) {
    return (
      <div
        className="flex items-center justify-center text-gray-500 text-xs"
        style={{ width, height }}
      >
        No data
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2">
      <svg width={width} height={height} className="overflow-visible">
        <path
          d={path}
          fill="none"
          stroke={strokeColor}
          strokeWidth="1.5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
      {showChange && (
        <span
          className={`text-xs font-medium ${
            isPositive ? 'text-green-400' : 'text-red-400'
          }`}
        >
          {isPositive ? '+' : ''}
          {changePercent.toFixed(1)}%
        </span>
      )}
    </div>
  );
}

// Simple seeded random number generator for deterministic sparklines
function seededRandom(seed: number): () => number {
  let state = seed;
  return () => {
    state = (state * 1103515245 + 12345) & 0x7fffffff;
    return (state / 0x7fffffff);
  };
}

// Generate fake sparkline data based on pool APY changes
// Uses seeded random so the same inputs always produce the same chart
export function generateSparklineData(
  currentApy: number,
  change24h: number,
  points: number = 12
): number[] {
  const data: number[] = [];
  const startApy = currentApy / (1 + change24h / 100);

  // Create a seed from the input values for deterministic output
  const seed = Math.floor((currentApy * 1000 + change24h * 100 + points) * 1000);
  const random = seededRandom(seed);

  for (let i = 0; i < points; i++) {
    const progress = i / (points - 1);
    const baseValue = startApy + (currentApy - startApy) * progress;
    const noise = (random() - 0.5) * currentApy * 0.1;
    data.push(Math.max(0, baseValue + noise));
  }

  // Ensure last point matches current APY
  data[data.length - 1] = currentApy;

  return data;
}

/**
 * Shared formatting utilities
 */

/**
 * Format a number with appropriate suffix (K, M, B)
 */
export function formatNumber(num: number | undefined, decimals = 2): string {
  if (num === undefined || num === null) return '-';
  if (num >= 1e9) return `$${(num / 1e9).toFixed(decimals)}B`;
  if (num >= 1e6) return `$${(num / 1e6).toFixed(decimals)}M`;
  if (num >= 1e3) return `$${(num / 1e3).toFixed(decimals)}K`;
  return `$${num.toFixed(decimals)}`;
}

/**
 * Format a percentage value
 */
export function formatPercent(value: number | undefined, decimals = 2): string {
  if (value === undefined || value === null) return '-';
  return `${value.toFixed(decimals)}%`;
}

/**
 * Format a date relative to now (e.g., "2h ago", "3d ago")
 */
export function formatTimeAgo(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) return `${days}d ago`;
  if (hours > 0) return `${hours}h ago`;
  if (minutes > 0) return `${minutes}m ago`;
  return 'Just now';
}

/**
 * Get color class based on score value
 */
export function getScoreColor(score: number): string {
  if (score >= 80) return 'text-green-400';
  if (score >= 60) return 'text-yellow-400';
  return 'text-red-400';
}

/**
 * Get color class based on change value (positive/negative)
 */
export function getChangeColor(change: number): string {
  if (change > 0) return 'text-green-400';
  if (change < 0) return 'text-red-400';
  return 'text-gray-400';
}

/**
 * Get badge color class based on opportunity type
 */
export function getTypeColor(type: string): string {
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
}

/**
 * Get badge color class based on risk level
 */
export function getRiskColor(risk: string): string {
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
}

/**
 * Get human-readable label for opportunity type
 */
export function getTypeLabel(type: string): string {
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
}

/**
 * Safely get the first element of an array or return fallback
 */
export function safeFirst<T>(arr: T[] | undefined | null, fallback: T): T {
  if (!arr || arr.length === 0) return fallback;
  return arr[0];
}

/**
 * Safely get lowercase of first array element or fallback
 */
export function safeFirstLower(arr: string[] | undefined | null, fallback: string): string {
  if (!arr || arr.length === 0) return fallback.toLowerCase();
  return arr[0].toLowerCase();
}

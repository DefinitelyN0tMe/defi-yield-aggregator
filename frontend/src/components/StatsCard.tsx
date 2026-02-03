import { ReactNode } from 'react';

interface StatsCardProps {
  title: string;
  value: ReactNode;
  subtitle?: string;
  icon?: ReactNode;
  change?: number;
  changeLabel?: string;
  variant?: 'default' | 'primary' | 'success' | 'warning';
}

export function StatsCard({
  title,
  value,
  subtitle,
  icon,
  change,
  changeLabel,
  variant = 'default',
}: StatsCardProps) {
  const variantStyles = {
    default: 'border-dark-700',
    primary: 'border-primary-500/30 bg-primary-500/5',
    success: 'border-green-500/30 bg-green-500/5',
    warning: 'border-yellow-500/30 bg-yellow-500/5',
  };

  const iconStyles = {
    default: 'bg-dark-700 text-gray-400',
    primary: 'bg-primary-500/20 text-primary-400',
    success: 'bg-green-500/20 text-green-400',
    warning: 'bg-yellow-500/20 text-yellow-400',
  };

  return (
    <div className={`card ${variantStyles[variant]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-sm font-medium text-gray-400">{title}</p>
          <p className="mt-2 text-2xl font-bold text-white">{value}</p>
          {subtitle && (
            <p className="mt-1 text-sm text-gray-500">{subtitle}</p>
          )}
          {change !== undefined && (
            <div className="mt-2 flex items-center gap-1">
              <span
                className={`text-sm font-medium ${
                  change >= 0 ? 'text-green-400' : 'text-red-400'
                }`}
              >
                {change >= 0 ? '+' : ''}
                {change.toFixed(2)}%
              </span>
              {changeLabel && (
                <span className="text-sm text-gray-500">{changeLabel}</span>
              )}
            </div>
          )}
        </div>
        {icon && (
          <div className={`p-3 rounded-lg ${iconStyles[variant]}`}>
            {icon}
          </div>
        )}
      </div>
    </div>
  );
}

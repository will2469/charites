import React, { useState } from 'react';

interface MetricCardProps {
  label: string;
  value: number;
  trend: 'up' | 'down' | 'neutral';
  loading?: boolean;
}

export const MetricCard: React.FC<MetricCardProps> = ({ label, value, trend, loading = false }) => {
  const [expanded, setExpanded] = useState<boolean>(false);
  const isPositive = trend === 'up';

  return (
    <>
      <div
        className={`relative overflow-hidden rounded-2xl p-6 transition-all ${
          loading ? 'opacity-50 pointer-events-none' : 'opacity-100'
        } bg-surface border border-slate-200 shadow-md hover:shadow-lg`}
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium text-slate-500">{label}</span>
          <span
            className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-semibold ${
              isPositive ? 'text-emerald-700 bg-emerald-50' : 'text-rose-700 bg-rose-50'
            }`}
          >
            {value < 0 ? `${value}%` : `+${value}%`}
          </span>
        </div>

        <div className="mt-4 flex items-baseline justify-between">
          <p className="text-2xl font-bold tracking-tight text-slate-900">{value.toLocaleString()}</p>
          <button
            type="button"
            onClick={() => setExpanded(!expanded)}
            className="text-xs font-medium text-primary hover:text-primary-hover"
            aria-expanded={expanded}
          >
            {expanded ? 'Hide details' : 'View breakdown'}
          </button>
        </div>

        <input type="hidden" name="metric-id" value={label} />
        {expanded && (
          <div className="mt-4 pt-4 border-t border-slate-100">
            <p className="text-xs text-slate-400">Detailed analytics breakdown will be loaded here.</p>
          </div>
        )}
      </div>
    </>
  );
};

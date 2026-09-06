import React from 'react';
import type { ChartConfiguration } from 'chart.js';

export interface ChartWrapperProps {
  config: ChartConfiguration;
}

export const ChartConfigSummary = ({ config }: ChartWrapperProps) => {
  return (
    <div>
      <span>Type: {config.type}</span>
    </div>
  );
};

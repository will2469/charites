import React from 'react';
import { Chart } from 'chart.js';

export const AnalyticsDashboard = ({ data }: { data: any }) => {
  return (
    <div>
      <h2>Analytics</h2>
      <Chart data={data} />
    </div>
  );
};

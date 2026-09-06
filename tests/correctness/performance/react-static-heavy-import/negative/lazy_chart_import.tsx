import React, { Suspense, lazy } from 'react';

const Chart = lazy(() => import('chart.js'));

export const AnalyticsDashboard = ({ data }: { data: any }) => {
  return (
    <div>
      <h2>Analytics</h2>
      <Suspense fallback={<div>Loading chart...</div>}>
        <Chart data={data} />
      </Suspense>
    </div>
  );
};

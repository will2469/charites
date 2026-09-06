import React from 'react';

export const DynamicBadge = ({ color }: { color: 'red' | 'blue' }) => {
  return <span className={`bg-${color}-100 text-${color}-800`}>Status</span>;
};

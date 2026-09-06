import React from 'react';

const COLOR_MAP = {
  red: 'bg-red-100 text-red-800',
  blue: 'bg-blue-100 text-blue-800',
} as const;

export const DynamicBadge = ({ color }: { color: 'red' | 'blue' }) => {
  return <span className={COLOR_MAP[color]}>Status</span>;
};

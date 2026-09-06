import React from 'react';

export const StatusBadge = ({ isSuccess }: { isSuccess: boolean }) => {
  return (
    <span className={`px-2 py-1 ${isSuccess ? 'bg-emerald-500 text-white' : 'bg-rose-500 text-white'}`}>
      Status
    </span>
  );
};

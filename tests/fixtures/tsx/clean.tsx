import React from 'react';

export const CleanCard: React.FC<{ title: string }> = ({ title }) => {
  return (
    <div className="bg-primary text-white p-4 rounded-md shadow-sm">
      <h2 className="text-lg font-semibold">{title}</h2>
      <p className="text-muted">Compliant component without token violations.</p>
    </div>
  );
};

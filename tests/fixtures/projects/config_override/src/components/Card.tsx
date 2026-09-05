import React from 'react';

export const WarningCard: React.FC = () => {
  return (
    <div className="bg-primary/10 p-4 rounded">
      <span>This violation is downgraded to warning via charites.yaml</span>
    </div>
  );
};

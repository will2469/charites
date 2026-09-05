import React from 'react';

export const SuppressedCard: React.FC = () => {
  return (
    <div>
      {/* charites:ignore theme.hardcode-opacity-color */}
      <div className="bg-primary/10">Suppressed inline</div>
    </div>
  );
};

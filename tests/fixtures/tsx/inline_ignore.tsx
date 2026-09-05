import React from 'react';

export const InlineIgnoreComponent: React.FC = () => {
  return (
    <div>
      {/* charites:ignore theme.hardcode-opacity-color */}
      <div className="bg-primary/10">Suppressed JSX block</div>
      {/* charites:ignore */}
      <span className="text-secondary/50">Suppressed all</span>
    </div>
  );
};

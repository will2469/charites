import React from 'react';

export const Card: React.FC = () => {
  return (
    <div className="bg-primary/10 text-secondary/50 p-6 rounded-lg border border-border">
      <h3 className="text-lg font-semibold">Violation Card</h3>
      <p className="text-muted">Card with hardcoded opacity classes.</p>
    </div>
  );
};

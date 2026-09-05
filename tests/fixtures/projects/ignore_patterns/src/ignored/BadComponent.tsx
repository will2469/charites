import React from 'react';

export const BadComponent: React.FC = () => {
  return <div className="bg-primary/10 text-secondary/30">Ignored by .charitesignore</div>;
};

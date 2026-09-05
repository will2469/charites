import React from 'react';

interface Props {
  variant: 'active' | 'inactive';
}

export const TemplateLiteralComponent: React.FC<Props> = ({ variant }) => {
  const dynamic = variant === 'active' ? 'ring-2' : '';
  return (
    <div className={`p-4 rounded-md ${dynamic} bg-primary/20`}>
      <span>Dynamic template string with opacity violation</span>
    </div>
  );
};

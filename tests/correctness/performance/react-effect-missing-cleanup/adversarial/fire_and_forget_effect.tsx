import React, { useEffect } from 'react';

export const PageTracker = ({ title }: { title: string }) => {
  useEffect(() => {
    document.title = title;
  }, [title]);

  return <div>Page: {title}</div>;
};

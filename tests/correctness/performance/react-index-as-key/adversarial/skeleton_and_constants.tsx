import React from 'react';

const SKELETON_SLOTS = 5;
const STATIC_NAV_TABS = ['Overview', 'Details', 'Settings'];

export const FeedSkeleton = () => {
  return (
    <div>
      <nav>
        {STATIC_NAV_TABS.map((tab, i) => (
          <span key={i}>{tab}</span>
        ))}
      </nav>
      <div className="skeleton-grid">
        {Array.from({ length: SKELETON_SLOTS }).map((_, i) => (
          <div key={i} className="skeleton-box" />
        ))}
      </div>
    </div>
  );
};

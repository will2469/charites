import React from 'react';

export function ThumbnailCard() {
  return (
    <div className="card p-4">
      {/* Small thumbnail without hero role and without high priority is completely valid */}
      <img src="/assets/thumb.png" alt="Thumbnail" width="48" height="48" className="w-12 h-12 rounded" />
    </div>
  );
}

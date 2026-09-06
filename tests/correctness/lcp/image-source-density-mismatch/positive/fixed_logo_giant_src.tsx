import React from 'react';

export function HeaderHero() {
  return (
    <header data-perf-role="hero">
      <img
        src="/assets/logo-2000.webp"
        width="120"
        height="40"
        alt="Corporate Logo"
        fetchpriority="high"
      />
    </header>
  );
}

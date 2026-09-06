import React from 'react';

export function HeaderHero() {
  return (
    <header data-perf-role="hero">
      <img
        src="/assets/logo-120.webp"
        srcSet="/assets/logo-120.webp 1x, /assets/logo-240.webp 2x"
        width="120"
        height="40"
        alt="Corporate Logo"
        fetchpriority="high"
      />
    </header>
  );
}

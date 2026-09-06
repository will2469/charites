import React from 'react';

export function HeroBanner() {
  return (
    <header className="hero-banner" data-perf-role="hero">
      <img src="/hero.webp" alt="Primary Banner" fetchpriority="high" className="w-full aspect-video" />
    </header>
  );
}

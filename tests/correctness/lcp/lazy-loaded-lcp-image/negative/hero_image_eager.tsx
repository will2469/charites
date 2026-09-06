import React from 'react';

export function HeroBanner() {
  return (
    <section className="hero-section" data-perf-role="hero">
      <h1>Welcome to Our Platform</h1>
      <img src="/assets/hero.webp" alt="Hero Banner" loading="eager" fetchpriority="high" className="w-full h-auto" />
    </section>
  );
}

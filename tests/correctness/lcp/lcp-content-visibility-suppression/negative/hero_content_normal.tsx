import React from 'react';

export function HeroBanner() {
  return (
    <section className="hero-section" data-perf-role="hero">
      <h1>Cloud Solution</h1>
      <img src="/hero.webp" fetchPriority="high" alt="Hero" />
    </section>
  );
}

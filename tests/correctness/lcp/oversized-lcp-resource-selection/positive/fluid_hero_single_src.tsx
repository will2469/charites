import React from 'react';

export function HeroBanner() {
  return (
    <section className="hero-section" data-perf-role="hero">
      <h1>Fluid Hero Banner</h1>
      <img
        src="/images/hero-3840x2160.jpg"
        alt="Hero Banner"
        className="w-full h-auto"
        fetchpriority="high"
      />
    </section>
  );
}

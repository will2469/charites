import React from 'react';

export function HeroBanner() {
  return (
    <section className="hero-section" data-perf-role="hero">
      <h1>Fluid Hero Banner</h1>
      <img
        src="/images/hero-1200.webp"
        srcset="/images/hero-400.webp 400w, /images/hero-800.webp 800w, /images/hero-1200.webp 1200w"
        sizes="100vw"
        alt="Hero Banner"
        className="w-full h-auto"
        fetchpriority="high"
      />
    </section>
  );
}

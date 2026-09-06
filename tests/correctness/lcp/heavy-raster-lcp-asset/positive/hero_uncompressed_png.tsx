import React from 'react';

export function HeroSection() {
  return (
    <section className="hero-container" data-perf-role="hero">
      <h1>Uncompressed Raster Hero</h1>
      <img
        src="/assets/landscape-hero.png"
        alt="Hero Landscape"
        fetchpriority="high"
        className="w-full h-auto"
      />
    </section>
  );
}

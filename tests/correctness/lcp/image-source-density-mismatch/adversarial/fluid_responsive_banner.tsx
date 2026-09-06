import React from 'react';

export function AdversarialMedia() {
  return (
    <div>
      {/* Fluid responsive hero: handled by oversized-lcp-resource-selection */}
      <section className="hero-banner" data-perf-role="hero">
        <img
          src="/assets/fluid-hero.webp"
          alt="Fluid Hero"
          className="w-full h-auto"
          fetchpriority="high"
        />
      </section>

      {/* SVG logo: vector media is density-independent */}
      <header data-perf-role="hero">
        <img
          src="/assets/logo.svg"
          width="120"
          height="40"
          alt="Vector Brand Logo"
        />
      </header>
    </div>
  );
}

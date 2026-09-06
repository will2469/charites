import React from 'react';

export function HeroAdversarial() {
  return (
    <section className="hero-container" data-perf-role="hero">
      <img
        src="/assets/vector-hero.svg"
        alt="Vector SVG"
        fetchpriority="high"
      />
      <img
        src="https://images.unsplash.com/photo-12345.png"
        alt="CDN Image with content negotiation"
        fetchpriority="high"
      />
      <img
        src="/assets/small-badge.png"
        width="32"
        height="32"
        alt="Small Badge"
        fetchpriority="high"
      />
      <picture>
        <source srcSet="/assets/hero.avif" type="image/avif" />
        <source srcSet="/assets/hero.webp" type="image/webp" />
        <img
          src="/assets/hero-fallback.png"
          alt="Hero with picture fallback"
          fetchpriority="high"
        />
      </picture>
    </section>
  );
}

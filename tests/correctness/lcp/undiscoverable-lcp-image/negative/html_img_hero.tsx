import React from 'react';

export function HeroHeader() {
  return (
    <header className="relative w-full h-[480px] overflow-hidden" data-perf-role="hero">
      <img src="/hero.webp" alt="Hero Background" fetchpriority="high" className="absolute inset-0 w-full h-full object-cover -z-10" />
      <h1 className="relative z-10 text-white p-8">Galactic Exploration</h1>
    </header>
  );
}

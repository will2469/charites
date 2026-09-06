import React from 'react';

export function HeroHeader() {
  return (
    <header className="w-full h-[480px] bg-[url('/hero.webp')] bg-cover" data-perf-role="hero">
      <h1 className="text-white">Galactic Exploration</h1>
    </header>
  );
}

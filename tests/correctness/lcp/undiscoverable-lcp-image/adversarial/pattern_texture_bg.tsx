import React from 'react';

export function TextureSection() {
  return (
    <header className="w-full h-48 bg-[url('/assets/grid.svg')] bg-repeat" data-perf-role="hero">
      {/* Decorative repeating pattern/texture is exempt from LCP image requirement */}
      <h1>Dashboard Header</h1>
    </header>
  );
}

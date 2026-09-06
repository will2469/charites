import React from 'react';

export function HeaderLogo() {
  return (
    <header data-perf-role="hero">
      <img
        src="/assets/logo.svg"
        alt="Brand Logo"
        fetchpriority="high"
      />
      <img
        src="/assets/fixed-avatar.webp"
        width="120"
        height="40"
        alt="Author Avatar"
        fetchpriority="high"
      />
    </header>
  );
}

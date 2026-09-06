import React from 'react';

export function HeaderWithNonRootLogo() {
  return (
    <header className="flex items-center justify-between px-6 py-4 border-b">
      <a href="/about" className="brand-logo flex items-center gap-2">
        <img src="/assets/brand.svg" alt="Brand Logo" className="h-8" />
        <span className="font-bold">Enterprise App</span>
      </a>
      <nav className="flex gap-4">
        <a href="/dashboard">Dashboard</a>
      </nav>
    </header>
  );
}

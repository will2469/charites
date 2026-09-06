import React from 'react';

export function LocalizedHeader() {
  return (
    <header className="flex items-center justify-between px-6 py-4 border-b">
      <a href="/id" aria-label="Beranda Portal" className="brand-logo flex items-center gap-2">
        <img src="/logo.svg" alt="Brand Logo" className="h-8" />
        <span className="font-bold">Portal</span>
      </a>
      <nav className="flex gap-4">
        <a href="/id/layanan">Layanan</a>
      </nav>
    </header>
  );
}

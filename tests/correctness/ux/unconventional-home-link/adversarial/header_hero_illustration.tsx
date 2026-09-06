import React from 'react';

export function HeaderWithIllustration() {
  return (
    <header className="flex items-center justify-between px-6 py-4">
      <a href="/" aria-label="Beranda">
        <img src="/logo.svg" alt="Company Logo" className="h-8" />
      </a>
      <div className="hero-decor flex items-center">
        <img src="/hero-banner-illustration.png" alt="Decorative Illustration" className="h-6" />
      </div>
    </header>
  );
}

import React from 'react';

export function FooterWidget() {
  return (
    <footer className="footer bg-muted">
      <p>Footer Content</p>
      {/* Below-fold images are allowed and encouraged to be lazy-loaded */}
      <img src="/assets/partner-logo.webp" alt="Partner Logo" loading="lazy" className="w-24 h-12" />
    </footer>
  );
}

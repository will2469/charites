import React from 'react';

export function UnreservedDynamicPage() {
  return (
    <main>
      <h1>Featured Articles</h1>
      <PromoBanner />
      <article>
        <p>Article content that gets pushed down when promo banner loads.</p>
      </article>
    </main>
  );
}

function PromoBanner() {
  return <div>Special Discount!</div>;
}

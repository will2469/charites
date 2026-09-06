import React from 'react';

export function AbsoluteDynamicPage() {
  return (
    <main>
      <h1>Featured Articles</h1>
      <PromoBanner className="fixed bottom-4 right-4 z-50 shadow-lg" />
      <NotificationWidget className="absolute top-2 right-2 z-40" />
      <article>
        <p>Article content remains unaffected by fixed overlays.</p>
      </article>
    </main>
  );
}

function PromoBanner({ className }: { className?: string }) {
  return <div className={className}>Floating Notice</div>;
}

function NotificationWidget({ className }: { className?: string }) {
  return <div className={className}>Absolute Badge</div>;
}

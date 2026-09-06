import React, { Suspense } from 'react';

export function ReservedDynamicPage() {
  return (
    <main>
      <h1>Featured Articles</h1>
      <div className="min-h-[120px]">
        <PromoBanner />
      </div>
      <Suspense fallback={<div className="h-24 bg-muted/20" />}>
        <NotificationWidget />
      </Suspense>
      <article>
        <p>Article content that stays stable.</p>
      </article>
    </main>
  );
}

function PromoBanner() {
  return <div>Special Discount!</div>;
}

function NotificationWidget() {
  return <div>Notification!</div>;
}

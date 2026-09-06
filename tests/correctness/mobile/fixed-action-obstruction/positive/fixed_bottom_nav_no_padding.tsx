import React from "react";

export function ObstructedLayout() {
  return (
    // Violation: fixed bottom nav without compensating padding on main or parent layout
    <div className="min-h-screen bg-background">
      <main className="p-4 space-y-4">
        <p>Konten formulir warga yang panjang...</p>
      </main>
      <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t">
        <button type="button">Beranda</button>
      </nav>
    </div>
  );
}

import React from "react";

export function ClearanceLayout() {
  return (
    // Compliant: main element provides pb-24 clearance
    <div className="min-h-screen bg-background">
      <main className="p-4 space-y-4 pb-24">
        <p>Konten formulir warga yang panjang...</p>
      </main>
      <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t pb-[env(safe-area-inset-bottom)]">
        <button type="button">Beranda</button>
      </nav>
    </div>
  );
}

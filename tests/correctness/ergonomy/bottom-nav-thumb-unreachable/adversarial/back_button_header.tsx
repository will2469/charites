import React from "react";

export function StandardHeaderPage({ onBack, onHelp }: { onBack: () => void; onHelp: () => void }) {
  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Header has navigation buttons (back and help) which must not be falsely flagged as unreachable primary CTAs */}
      <header className="sticky top-0 z-10 flex items-center justify-between p-4 bg-card border-b">
        <button type="button" onClick={onBack} aria-label="Kembali" className="p-2">
          <span>&lt;</span>
        </button>
        <h1 className="font-semibold text-lg">Pusat Bantuan Desa</h1>
        <button type="button" onClick={onHelp} aria-label="Bantuan Informasi" className="p-2">
          <span>?</span>
        </button>
      </header>
      <main className="p-4">
        <p>Konten panduan warga</p>
      </main>
    </div>
  );
}

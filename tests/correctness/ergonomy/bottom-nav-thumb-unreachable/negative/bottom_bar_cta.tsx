import React from "react";

export function CompliantProfileForm({ onBack, onSubmit }: { onBack: () => void; onSubmit: () => void }) {
  return (
    <div className="min-h-screen bg-background text-foreground pb-24">
      <header className="sticky top-0 z-10 flex items-center justify-between p-4 bg-card border-b">
        <button type="button" onClick={onBack} aria-label="Kembali" className="p-2">
          <span>Back</span>
        </button>
        <h1 className="font-semibold text-lg">Edit Profil Warga</h1>
      </header>
      <main className="p-4 space-y-4">
        <input name="nama" placeholder="Nama Lengkap" className="w-full h-11 px-3 border rounded-xl" />
      </main>
      <footer className="fixed bottom-0 inset-x-0 p-4 bg-card border-t pb-[env(safe-area-inset-bottom)]">
        <button
          type="submit"
          onClick={onSubmit}
          className="w-full h-12 bg-primary text-primary-foreground rounded-xl font-semibold active:scale-95"
        >
          Simpan Perubahan
        </button>
      </footer>
    </div>
  );
}

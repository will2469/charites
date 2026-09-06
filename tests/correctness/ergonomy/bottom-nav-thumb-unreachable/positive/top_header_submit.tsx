import React from "react";

export function EditProfilePage({ onBack, onSubmit }: { onBack: () => void; onSubmit: () => void }) {
  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Top sticky header has the only primary submit CTA - unreachable on tall mobile displays */}
      <header className="sticky top-0 z-10 flex items-center justify-between p-4 bg-card border-b">
        <button type="button" onClick={onBack} aria-label="Kembali" className="p-2">
          <span>Back</span>
        </button>
        <h1 className="font-semibold text-lg">Edit Profil Warga</h1>
        <button
          type="submit"
          onClick={onSubmit}
          className="h-10 px-4 bg-primary text-primary-foreground rounded-xl font-medium"
        >
          Simpan
        </button>
      </header>
      <main className="p-4 space-y-4">
        <input name="nama" placeholder="Nama Lengkap" className="w-full h-11 px-3 border rounded-xl" />
        <input name="nik" placeholder="NIK Warga" className="w-full h-11 px-3 border rounded-xl" />
      </main>
    </div>
  );
}

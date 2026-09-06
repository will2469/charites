import React from "react";

export function LockedDialog() {
  return (
    // Violation: modal dialog with overflow-hidden without internal scrollable region
    <div role="dialog" aria-modal="true" className="fixed inset-0 overflow-hidden flex items-center justify-center p-4">
      <div className="bg-card p-6 rounded-2xl w-full max-w-md h-screen">
        <h2>Form Permohonan Bantuan</h2>
        <div className="space-y-4">
          <input name="nama" />
        </div>
        <button type="submit">Kirim</button>
      </div>
    </div>
  );
}

import React from "react";

export function ScrollableDialog() {
  return (
    // Compliant: overflow-y-auto on dialog wrapper allows smooth scrolling
    <div role="dialog" aria-modal="true" className="fixed inset-0 overflow-y-auto flex items-center justify-center p-4">
      <div className="bg-card p-6 rounded-2xl w-full max-w-md my-auto">
        <h2>Form Permohonan Bantuan</h2>
        <div className="space-y-4">
          <input name="nama" />
        </div>
        <button type="submit">Kirim</button>
      </div>
    </div>
  );
}

import React from "react";

export function CompliantAriaDialog() {
  return (
    <div>
      {/* With aria-modal and aria-labelledby */}
      <div role="dialog" aria-modal="true" aria-labelledby="dialog-title" className="fixed inset-0 z-50">
        <h2 id="dialog-title">Judul Dialog</h2>
        <p>Konten</p>
      </div>

      {/* With aria-modal and aria-label */}
      <div role="alertdialog" aria-modal="true" aria-label="Konfirmasi Pembatalan" className="fixed inset-0 z-50">
        <p>Apakah yakin?</p>
      </div>
    </div>
  );
}

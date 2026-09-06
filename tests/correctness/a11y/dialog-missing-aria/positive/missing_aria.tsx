import React from "react";

export function MissingAriaModal() {
  return (
    <div>
      {/* Missing both aria-modal and aria-labelledby */}
      <div role="dialog" className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
        <div className="bg-background p-6 rounded-xl">
          <h2>Judul Dialog</h2>
          <p>Konten dialog</p>
        </div>
      </div>

      {/* Missing accessible name */}
      <div role="alertdialog" aria-modal="true" className="fixed inset-0 z-50">
        <p>Peringatan keras</p>
      </div>
    </div>
  );
}

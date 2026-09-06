import React from "react";

export function PWADocs() {
  // Catatan: jangan gunakan navigator.standalone langsung karena iOS-only
  // Gunakan matchMedia('(display-mode: standalone)')
  return (
    <div className="p-4">
      <p>Petunjuk instalasi aplikasi Progressive Web App.</p>
    </div>
  );
}

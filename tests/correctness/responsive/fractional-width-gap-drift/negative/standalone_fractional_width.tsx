import React from "react";

export function StandaloneFractionalWidth() {
  return (
    // Standalone fractional widths are not inherently defects without conflicting container constraints
    <div>
      <div className="w-1/2">Kolom 1</div>
      <div className="w-2/3">Konten Utama</div>
      <div className="w-full md:w-2/3">Responsive Mobile First</div>
    </div>
  );
}

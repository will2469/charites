import React from "react";

export function CssGridAlternative() {
  return (
    // CSS Grid automatically deducts gap space from fraction tracks
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="md:col-span-2">Kolom 1</div>
      <div>Kolom 2</div>
    </div>
  );
}

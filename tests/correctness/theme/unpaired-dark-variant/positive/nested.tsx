import React from "react";

// P2: Inverted container background with un-adapted child text colors
export function InvertedCard() {
  return (
    <div className="bg-white dark:bg-zinc-900 p-6">
      <span className="text-zinc-900">Judul Komponen</span>
      <p className="text-slate-800">Deskripsi tanpa dark:text</p>
    </div>
  );
}

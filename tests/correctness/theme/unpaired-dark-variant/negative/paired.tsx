import React from "react";

// N2: Symmetrically paired background and text variants
export function PairedCard() {
  return (
    <div className="bg-white dark:bg-zinc-900 p-6">
      <span className="text-zinc-900 dark:text-zinc-100">Paired Text</span>
    </div>
  );
}

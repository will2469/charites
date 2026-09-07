import React from "react";

export function DefaultNowrapHalfHalf() {
  return (
    // Default nowrap with default shrink-1 absorbs the 16px gap safely
    <div className="flex gap-4">
      <div className="w-1/2">Kolom 1</div>
      <div className="w-1/2">Kolom 2</div>
    </div>
  );
}

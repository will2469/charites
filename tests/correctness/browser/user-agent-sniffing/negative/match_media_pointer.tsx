import React from "react";

export function TouchCapabilityBanner() {
  const isTouchDevice = typeof window !== "undefined" && window.matchMedia("(pointer: coarse)").matches;

  return (
    <div className="p-4 bg-muted text-sm">
      {isTouchDevice ? <span>Sentuhan diaktifkan</span> : <span>Navigasi mouse/pointer presisi</span>}
    </div>
  );
}

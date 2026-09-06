import React from "react";

export function StandardPWABadge() {
  // Compliant: standard display-mode media query with iOS legacy fallback
  const isInstalled =
    (typeof window !== "undefined" && window.matchMedia("(display-mode: standalone)").matches) ||
    (typeof navigator !== "undefined" && Boolean((navigator as any).standalone));

  return (
    <div className="p-2 border">
      {isInstalled ? <span>Mode PWA Terpasang</span> : <span>Mode Browser Standar</span>}
    </div>
  );
}

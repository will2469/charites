import React from "react";

export function PWABadge() {
  // Direct iOS-only property: undefined on Android and Windows
  const isInstalled = navigator.standalone;

  return (
    <div className="p-2 border">
      {isInstalled ? <span>Mode PWA Terpasang</span> : <span>Mode Browser Standar</span>}
    </div>
  );
}

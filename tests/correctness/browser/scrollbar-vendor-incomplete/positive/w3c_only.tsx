import React from "react";

// P2: W3C standard scrollbar in inline style without WebKit scrollbar pairing
export function W3CScrollbox() {
  return (
    <div
      style={{ scrollbarWidth: "thin", scrollbarColor: "#cbd5e1 #f1f5f9" }}
      className="overflow-y-auto max-h-60"
    >
      <p>Scrollable items...</p>
    </div>
  );
}

import React from "react";

// A2: Splitting URL pathname is completely valid
export function Breadcrumbs() {
  const segments = typeof window !== "undefined" ? window.location.pathname.split('/') : [];
  return (
    <nav>
      {segments.map((seg, i) => (
        <span key={i}>{seg}</span>
      ))}
    </nav>
  );
}

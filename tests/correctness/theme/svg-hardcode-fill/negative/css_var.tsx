import React from "react";

// N2: Dynamic CSS variable in gradient stop
export function SemanticGradient() {
  return (
    <svg viewBox="0 0 100 100">
      <defs>
        <linearGradient id="grad2">
          <stop stop-color="var(--primary)" offset="0%" />
          <stop stop-color="var(--accent)" offset="100%" />
        </linearGradient>
      </defs>
    </svg>
  );
}

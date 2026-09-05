import React from "react";

// P2: Hardcoded stop-color in gradient stop
export function GradientDef() {
  return (
    <svg viewBox="0 0 100 100">
      <defs>
        <linearGradient id="grad">
          <stop stop-color="#3b82f6" offset="100%" />
        </linearGradient>
      </defs>
    </svg>
  );
}

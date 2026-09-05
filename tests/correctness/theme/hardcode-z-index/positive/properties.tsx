import React from "react";

export function ZIndexPropertyViolations() {
  return (
    <nav className="[z-index:1000]">
      <span className="z-[99]">Custom z layout</span>
    </nav>
  );
}

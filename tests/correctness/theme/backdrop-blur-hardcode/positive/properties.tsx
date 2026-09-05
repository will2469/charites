import React from "react";

export function BackdropBlurPropertyViolations() {
  return (
    <div className="[backdrop-filter:blur(7px)]">
      <span className="[filter:blur(4px)]">Blurred span</span>
    </div>
  );
}

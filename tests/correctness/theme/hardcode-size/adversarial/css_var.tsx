import React from "react";

export function CSSVariableArbitrarySizing() {
  return (
    <div className="p-[var(--spacing-custom)] w-[var(--container-max)] top-[var(--header-height)]">
      <span>Dynamic token variable sizes</span>
    </div>
  );
}

import React from "react";

export function CSSVariableArbitrarySizing() {
  return (
    <div className="p-[var(--spacing-custom)] w-[var(--container-max)] top-[var(--header-height)]">
      <div className="scale-[var(--scale-press)]">Dynamic token variable scale</div>
      <div className="[scale:var(--scale-hover)]">Dynamic token property scale</div>
    </div>
  );
}

import React from "react";

export function GapXZeroOverride() {
  return (
    // gap-x-0 overrides shorthand gap-4 for horizontal direction
    <div className="flex flex-wrap gap-4 gap-x-0">
      <div className="w-1/2">A</div>
      <div className="w-1/2">B</div>
    </div>
  );
}

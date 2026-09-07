import React from "react";

export function GapYOnlyNoDrift() {
  return (
    // gap-y-4 is vertical only, zero horizontal spacing drift
    <div className="flex flex-wrap gap-y-4">
      <div className="w-1/2">A</div>
      <div className="w-1/2">B</div>
    </div>
  );
}

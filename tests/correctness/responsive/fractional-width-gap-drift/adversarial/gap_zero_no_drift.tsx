import React from "react";

export function GapZeroNoDrift() {
  return (
    // gap-0 means zero spacing drift
    <div className="flex flex-wrap gap-0">
      <div className="w-1/2">A</div>
      <div className="w-1/2">B</div>
    </div>
  );
}

import React from "react";

export function NonFractionalShrink0() {
  return (
    // Non-fractional child with shrink-0 does not cause fractional allocation conflict
    <div className="flex gap-4">
      <div className="w-1/2">A</div>
      <div className="w-1/2">B</div>
      <div className="w-auto shrink-0">C</div>
    </div>
  );
}

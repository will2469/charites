import React from "react";

// Positive fixture: Arbitrary inline transform scale modifiers
export function ArbitraryScaleContainer() {
  return (
    <div className="space-y-4">
      <div className="active:scale-[0.99]">Active scale</div>
      <div className="-scale-[0.99]">Negative uniform scale</div>
      <div className="hover:scale-[1.02]">Hover scale</div>
      <div className="-scale-x-[0.95]">Negative scale x</div>
      <div className="scale-y-[1.05]">Scale y</div>
      <div className="scale-z-[1.05]">Scale z</div>
      <div className="[scale:0.98]">Arbitrary property scale</div>
    </div>
  );
}

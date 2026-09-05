import React from "react";

export function FocusRingVariantViolations() {
  return (
    <div>
      <button className="focus-visible:ring-[#2563eb]">Action</button>
      <input className="outline-[#3b82f6]" />
      <span className="[ring-color:#ff0000]">Highlighted</span>
    </div>
  );
}

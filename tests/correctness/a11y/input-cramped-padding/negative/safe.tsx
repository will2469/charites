import React from "react";

export function CompliantPaddingInputs() {
  return (
    <div>
      <input className="h-11 px-3.5 py-2.5 border" placeholder="Name" />
      <input className="h-8 min-h-[44px] px-3" placeholder="Compensated" />
      <textarea className="py-3 px-4 border" placeholder="Bio" />
    </div>
  );
}

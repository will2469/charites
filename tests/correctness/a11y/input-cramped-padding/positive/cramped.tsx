import React from "react";

export function CrampedPaddingInputs() {
  return (
    <div>
      <input className="h-8 px-2 text-sm border" placeholder="Cramped height" />
      <textarea className="py-1 px-2 border" placeholder="Cramped padding" />
    </div>
  );
}

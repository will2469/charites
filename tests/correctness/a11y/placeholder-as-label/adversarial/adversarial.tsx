import React from "react";

export function AdversarialLabelInputs() {
  return (
    <div>
      {/* Title attribute provides valid accessible name */}
      <input id="q" placeholder="Cari..." title="Pencarian" />
      {/* Inputs without placeholders are not flagged by this rule */}
      <input type="text" id="simple-input" />
      {/* Non-text inputs */}
      <input type="checkbox" />
      <input type="submit" value="Kirim" />
    </div>
  );
}

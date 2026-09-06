import React from "react";

export function AdversarialZoomInputs() {
  return (
    <div>
      {/* Checkboxes and radios do not take text entry and do not trigger Safari auto-zoom */}
      <input type="checkbox" className="text-sm rounded" />
      <input type="radio" className="text-xs" />
      <input type="hidden" className="text-sm" />
      {/* Non-input elements with text-sm */}
      <p className="text-sm">Description text</p>
      <span className="text-xs">Badge</span>
    </div>
  );
}

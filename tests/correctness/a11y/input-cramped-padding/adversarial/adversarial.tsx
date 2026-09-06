import React from "react";

export function AdversarialPaddingInputs() {
  return (
    <div>
      {/* Checkbox has naturally small h-4 w-4 dimensions and is not a textual input */}
      <input type="checkbox" className="h-4 w-4" />
      <input type="radio" className="h-4 w-4" />
      {/* Non-input div with cramped padding is not a form input */}
      <div className="h-8 py-1">status tag</div>
    </div>
  );
}

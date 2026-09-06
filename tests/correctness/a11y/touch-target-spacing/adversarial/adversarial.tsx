import React from "react";

export function AdversarialSpacing() {
  return (
    <div>
      {/* Container with margin-compensated children */}
      <div className="flex">
        <button className="h-11 w-11 mx-2">Item 1</button>
        <button className="h-11 w-11 mx-2">Item 2</button>
      </div>
      {/* Non-interactive items in flex */}
      <div className="flex gap-0">
        <span className="text-sm">Label 1:</span>
        <span className="text-sm">Value 1</span>
      </div>
    </div>
  );
}

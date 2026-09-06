import React from "react";

// N2: Plain scrollable container without custom scrollbar styling
export function PlainScrollContainer() {
  return (
    <div className="overflow-y-auto max-h-96 p-4 border rounded-lg">
      <ul className="space-y-2">
        <li>Item 1</li>
        <li>Item 2</li>
        <li>Item 3</li>
      </ul>
    </div>
  );
}

import React from "react";

// N2: Standard button with color transition on hover - not an element reveal
export function ColorHoverButton() {
  return (
    <button className="bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded-lg">
      Click Me
    </button>
  );
}

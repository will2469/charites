import React from "react";

// N2: Custom styled checkbox with explicit appearance-none
export function CheckedReset() {
  return (
    <div className="flex items-center space-x-2">
      <input
        type="checkbox"
        className="appearance-none h-4 w-4 rounded border-gray-300 bg-white checked:bg-blue-600 focus:ring-2"
      />
      <span>Terms and conditions</span>
    </div>
  );
}

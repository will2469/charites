import React from "react";

// P2: Custom styled checkbox without appearance-none
export function CustomCheckbox() {
  return (
    <div className="flex items-center space-x-2">
      <input
        type="checkbox"
        className="h-4 w-4 rounded border-gray-300 bg-white text-blue-600 focus:ring-2"
      />
      <span>Remember me</span>
    </div>
  );
}

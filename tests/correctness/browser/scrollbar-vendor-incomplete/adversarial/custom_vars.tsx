import React from "react";

// A2: CSS variable names or class names that do not configure scrollbar properties
export function CustomScrollbarVarHolder() {
  return (
    <div
      className="p-4"
      style={{
        ["--custom-scrollbar-width" as any]: "12px",
        ["--custom-scrollbar-color" as any]: "blue",
      }}
    >
      <span className="text-sm text-gray-500">
        Variable holder without scrollbar declarations
      </span>
    </div>
  );
}

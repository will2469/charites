import React from "react";

// A1: Custom headless Select component and unstyled native controls
interface SelectProps {
  className?: string;
  children: React.ReactNode;
}

export function Select({ className, children }: SelectProps) {
  return (
    <div role="combobox" aria-expanded="false" className={`border rounded bg-white ${className || ""}`}>
      {children}
    </div>
  );
}

export function CustomSelectWrapper() {
  return (
    <Select className="p-2 border-gray-300">
      <span>Custom Dropdown Trigger</span>
    </Select>
  );
}

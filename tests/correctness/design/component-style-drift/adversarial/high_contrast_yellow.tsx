import React from "react";

// A3: High contrast yellow background with dark text - passes WCAG 2.2 AA (> 13:1)
export function HighContrastWarning() {
  return (
    <Button className="bg-yellow-400 text-black rounded-md font-semibold">
      Accessible Warning
    </Button>
  );
}

import React from "react";

// A4: Hover and focus variant modifiers on background and text - must not cause false positive in static evaluation
export function VariantButton() {
  return (
    <Button className="bg-black text-white hover:bg-yellow-400 hover:text-black focus-visible:ring-2">
      Interactive Button
    </Button>
  );
}

import React from "react";

export function FlexOneFluid() {
  return (
    // flex-1 and flex-1 basis-0 distribute space fluidly
    <div className="flex gap-4">
      <div className="flex-1 basis-0">Item A</div>
      <div className="flex-1 basis-0">Item B</div>
    </div>
  );
}

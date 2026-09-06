import React from "react";

export function SizingViolations() {
  return (
    <section className="h-[450px] min-w-[280px] max-h-[600px] top-[14px]">
      <span className="leading-[23px] tracking-[0.7px]">Text sizing</span>
      <div className="p-3.25 w-2.75">Non-standard fractional scale</div>
    </section>
  );
}

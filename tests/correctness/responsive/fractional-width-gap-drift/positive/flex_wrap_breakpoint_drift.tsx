import React from "react";

export function FlexWrapBreakpointDrift() {
  return (
    <div className="flex flex-col md:flex-row md:flex-wrap md:gap-4">
      <div className="w-full md:w-2/3">Konten Utama</div>
      <div className="w-full md:w-1/3">Sidebar</div>
    </div>
  );
}

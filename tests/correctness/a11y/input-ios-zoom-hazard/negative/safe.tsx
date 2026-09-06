import React from "react";

export function CompliantZoomInputs() {
  return (
    <div>
      {/* 16px mobile baseline with responsive sm:text-sm */}
      <input className="text-base sm:text-sm px-3.5 py-2.5 border rounded" placeholder="Email" />
      {/* Standard 16px base font */}
      <input className="text-base border rounded px-3 py-2" placeholder="Name" />
      {/* sm:text-sm alone does not set mobile font, preserving default 16px on mobile */}
      <input className="sm:text-sm border rounded px-3 py-2" placeholder="Phone" />
    </div>
  );
}

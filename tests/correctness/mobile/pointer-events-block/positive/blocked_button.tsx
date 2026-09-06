import React from "react";

export function BlockedButton() {
  return (
    // Violation: interactive button inside pointer-events-none parent without pointer-events-auto
    <div className="pointer-events-none opacity-90 p-4">
      <button onClick={() => console.log("clicked")} className="bg-primary text-white px-4 py-2">
        Simpan Data
      </button>
    </div>
  );
}

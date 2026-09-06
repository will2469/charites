import React from "react";

export function RestoredButton() {
  return (
    // Compliant: pointer-events-auto restores touch event hit testing
    <div className="pointer-events-none opacity-90 p-4">
      <button onClick={() => console.log("clicked")} className="pointer-events-auto bg-primary text-white px-4 py-2 rounded-xl">
        Simpan Data
      </button>
    </div>
  );
}

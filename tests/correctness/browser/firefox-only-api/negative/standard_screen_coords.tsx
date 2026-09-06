import React from "react";

export function UniversalScreenTracker() {
  const getCoordinates = () => {
    // Compliant: standard screenX
    const x = typeof window !== "undefined" ? window.screenX : 0;
    return x;
  };

  return (
    <div className="p-4">
      <button type="button" onClick={getCoordinates}>
        Ambil Posisi Layar Standar
      </button>
    </div>
  );
}

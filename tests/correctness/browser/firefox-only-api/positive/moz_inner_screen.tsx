import React from "react";

export function ScreenTracker() {
  const getCoordinates = () => {
    // Non-standard Gecko extension
    const x = window.mozInnerScreenX;
    return x;
  };

  return (
    <div className="p-4">
      <button type="button" onClick={getCoordinates}>
        Ambil Posisi Layar
      </button>
    </div>
  );
}

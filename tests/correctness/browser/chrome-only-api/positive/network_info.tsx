import React from "react";

export function QualityIndicator() {
  // Unguarded: crashes on Safari & desktop Firefox
  const speed = navigator.connection.effectiveType;

  return (
    <div className="p-2 border">
      <span>Kecepatan jaringan: {speed}</span>
    </div>
  );
}

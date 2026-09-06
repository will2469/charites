import React from "react";

export function GuardedQualityIndicator() {
  // Compliant: optional chaining with fallback
  const speed = navigator.connection?.effectiveType || "4g";

  return (
    <div className="p-2 border">
      <span>Kecepatan jaringan: {speed}</span>
    </div>
  );
}

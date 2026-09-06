import React from "react";

// N2: Fully valid -webkit-line-clamp triad with display: -webkit-box and -webkit-box-orient: vertical
export function ValidClamp() {
  return (
    <div
      style={{
        display: "-webkit-box",
        WebkitBoxOrient: "vertical",
        WebkitLineClamp: 2,
        overflow: "hidden",
      }}
      className="text-gray-900"
    >
      <p>Correctly clamped two line heading.</p>
    </div>
  );
}

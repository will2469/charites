import React from "react";

// P2: Incomplete line-clamp triad missing display: -webkit-box and -webkit-box-orient: vertical
export function BrokenClamp() {
  return (
    <div
      style={{
        WebkitLineClamp: 3,
        overflow: "hidden",
      }}
      className="text-gray-700"
    >
      <p>Long description that will fail to truncate on Safari and Chromium without the box triad.</p>
    </div>
  );
}

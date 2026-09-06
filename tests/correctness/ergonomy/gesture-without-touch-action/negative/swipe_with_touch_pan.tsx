import React from "react";

export function CompliantSwipeableGallery() {
  const handleTouchStart = () => {};
  const handleTouchMove = () => {};

  return (
    // Compliant: touch-pan-y isolates horizontal swipe and preserves vertical scroll
    <div
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      className="flex overflow-x-auto gap-4 p-4 touch-pan-y"
    >
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 1</div>
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 2</div>
    </div>
  );
}

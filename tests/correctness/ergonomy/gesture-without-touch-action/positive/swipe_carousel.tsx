import React from "react";

export function SwipeableCardGallery() {
  const handleTouchStart = () => {};
  const handleTouchMove = () => {};

  return (
    // Custom touchmove swipe handler without touch-action causes compositor scroll clashing
    <div
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      className="flex overflow-x-auto gap-4 p-4"
    >
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 1</div>
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 2</div>
    </div>
  );
}

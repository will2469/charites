import React from "react";

export function ClickOnlyCarousel() {
  const handleClick = () => {};

  return (
    // Click only: no swipe/drag event listener, should NOT require touch-action
    <div onClick={handleClick} className="flex gap-4 p-4">
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 1</div>
      <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu Layanan 2</div>
    </div>
  );
}

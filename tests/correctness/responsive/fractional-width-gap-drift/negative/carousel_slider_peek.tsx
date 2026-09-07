import React from "react";

export function CarouselSliderPeek() {
  return (
    // Intent-driven carousel with overflow-x-auto and shrink-0 peek cards
    <div className="flex overflow-x-auto gap-4">
      <div className="w-4/5 shrink-0">Slide 1</div>
      <div className="w-4/5 shrink-0">Slide 2</div>
    </div>
  );
}

import React from 'react';

export function SwipeableWidget({ handleDragStart, handleDragMove }: any) {
  return (
    <div onPointerDown={handleDragStart} onPointerMove={handleDragMove} className="w-full h-48 bg-muted touch-none">
      <span>Touch None Drag Handle</span>
    </div>
  );
}

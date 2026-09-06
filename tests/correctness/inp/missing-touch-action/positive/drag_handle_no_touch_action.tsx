import React from 'react';

export function CustomDragHandle({ handleDragStart, handleDragMove }: any) {
  return (
    <div onPointerDown={handleDragStart} onPointerMove={handleDragMove} className="w-full h-48 bg-muted">
      <span>Drag Handle</span>
    </div>
  );
}

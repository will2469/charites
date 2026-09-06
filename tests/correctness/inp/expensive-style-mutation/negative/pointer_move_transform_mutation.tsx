import React from 'react';

export function HardwareTransformCard() {
  return (
    <div
      onPointerMove={(e) => {
        e.currentTarget.style.transform = `translateY(${e.clientY / 10}px)`;
      }}
    >
      Interactive Card
    </div>
  );
}

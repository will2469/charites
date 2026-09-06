import React from 'react';

export function InteractiveCard() {
  return (
    <div
      onPointerMove={(e) => {
        e.currentTarget.style.boxShadow = `0 ${e.clientY / 10}px 30px rgba(0,0,0,0.5)`;
        e.currentTarget.style.filter = `blur(${e.clientX / 50}px)`;
      }}
    >
      Interactive Card
    </div>
  );
}

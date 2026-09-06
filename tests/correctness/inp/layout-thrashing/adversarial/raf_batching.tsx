import React from 'react';

export function RafBoxResizer() {
  return (
    <button
      onClick={() => {
        const el = document.getElementById('box')!;
        el.style.width = '200px';
        requestAnimationFrame(() => {
          const height = el.offsetHeight;
          el.style.height = `${height * 2}px`;
        });
      }}
    >
      Resize Box with RAF
    </button>
  );
}

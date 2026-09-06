import React from 'react';

export function BoxResizer() {
  return (
    <button
      onClick={() => {
        const el = document.getElementById('box')!;
        const height = el.offsetHeight;
        el.style.width = '200px';
        el.style.height = `${height * 2}px`;
      }}
    >
      Resize Box Safely
    </button>
  );
}

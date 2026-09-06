import React from 'react';

export function BoxResizer() {
  return (
    <button
      onClick={() => {
        const el = document.getElementById('box')!;
        el.style.width = '200px';
        const height = el.offsetHeight;
        el.style.height = `${height * 2}px`;
      }}
    >
      Resize Box
    </button>
  );
}

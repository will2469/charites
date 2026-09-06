import React from 'react';

export function ShortTabLoop() {
  const tabs = ['one', 'two', 'three'];
  return (
    <button
      onClick={() => {
        const names: string[] = [];
        for (let i = 0; i < 3; i++) {
          names.push(tabs[i]);
        }
      }}
    >
      Iterate Fixed
    </button>
  );
}

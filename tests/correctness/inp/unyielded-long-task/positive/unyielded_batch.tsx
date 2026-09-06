import React from 'react';

export function UnyieldedBatchProcessor({ items }: { items: string[] }) {
  return (
    <button
      onClick={() => {
        for (let i = 0; i < items.length; i++) {
          heavyCalculation(items[i]);
        }
      }}
    >
      Process Batch
    </button>
  );
}

function heavyCalculation(item: string) {
  // compute
}

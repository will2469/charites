import React from 'react';

export function CooperativeBatchProcessor({ items }: { items: string[] }) {
  return (
    <button
      onClick={async () => {
        for (let i = 0; i < items.length; i++) {
          heavyCalculation(items[i]);
          if (i % 50 === 0) {
            await (window as any).scheduler?.yield?.();
          }
        }
      }}
    >
      Process Cooperatively
    </button>
  );
}

function heavyCalculation(item: string) {
  // compute
}

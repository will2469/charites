import React from 'react';

export function YieldedDataViewer({ rawPayload }: { rawPayload: string }) {
  return (
    <button
      onClick={async () => {
        await (window as any).scheduler?.yield?.();
        const items = JSON.parse(rawPayload);
        const sorted = items.sort((a: any, b: any) => b.score - a.score);
        console.log(sorted);
      }}
    >
      Sort Large Payload Cooperatively
    </button>
  );
}

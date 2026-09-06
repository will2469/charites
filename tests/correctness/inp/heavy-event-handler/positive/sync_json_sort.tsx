import React from 'react';

export function HeavyDataViewer({ rawPayload }: { rawPayload: string }) {
  return (
    <button
      onClick={() => {
        const items = JSON.parse(rawPayload);
        const sorted = items.sort((a: any, b: any) => b.score - a.score);
        console.log(sorted);
      }}
    >
      Sort Large Payload
    </button>
  );
}

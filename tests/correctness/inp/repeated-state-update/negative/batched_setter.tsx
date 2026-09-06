import React, { useState } from 'react';

export function BatchedAsyncUpdater({ items }: { items: string[] }) {
  const [details, setDetails] = useState<any[]>([]);

  return (
    <button
      onClick={async () => {
        const results = [];
        for (const item of items) {
          const res = await fetch(`/api/detail/${item}`);
          results.push(await res.json());
        }
        setDetails((prev) => [...prev, ...results]);
      }}
    >
      Fetch Details Batched
    </button>
  );
}

import React, { useState } from 'react';

export function AsyncLoopUpdater({ items }: { items: string[] }) {
  const [details, setDetails] = useState<any[]>([]);

  return (
    <button
      onClick={async () => {
        for (const item of items) {
          const res = await fetch(`/api/detail/${item}`);
          const json = await res.json();
          setDetails((prev) => [...prev, json]);
        }
      }}
    >
      Fetch Details Sequentially
    </button>
  );
}

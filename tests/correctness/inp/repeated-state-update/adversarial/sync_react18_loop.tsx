import React, { useState } from 'react';

export function SyncLoopComponent() {
  const [total, setTotal] = useState(0);

  return (
    <button
      onClick={() => {
        let sum = 0;
        for (let i = 0; i < 10; i++) {
          sum += i;
        }
        setTotal(sum);
      }}
    >
      Calculate Sum
    </button>
  );
}

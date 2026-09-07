import React from 'react';

export function NestedDifferentParents() {
  return (
    <div className="flex flex-col gap-2">
      <p>
        Item pertama baris 1<br />Item pertama baris 2
      </p>
      <p>
        Item kedua baris 1<br />Item kedua baris 2
      </p>
    </div>
  );
}

import React from 'react';

export function NativeDetailsDisclosure() {
  return (
    <details className="border rounded p-4">
      <summary className="font-semibold cursor-pointer">View Details</summary>
      <div className="mt-2">
        <p>Native browser disclosure widget without max-height hacks.</p>
      </div>
    </details>
  );
}

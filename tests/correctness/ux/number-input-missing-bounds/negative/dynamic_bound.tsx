import React from "react";

export function DynamicBoundInput({ domainMin }: { domainMin: number }) {
  return (
    <input
      type="number"
      name="batch_size"
      min={domainMin}
      onWheel={(e) => e.currentTarget.blur()}
    />
  );
}

import React from "react";

export function BoundedNegative() {
  return (
    <input
      type="number"
      name="temperature"
      min="-50"
      onWheel={(e) => e.currentTarget.blur()}
    />
  );
}

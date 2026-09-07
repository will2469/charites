import React from "react";

export function BoundedZero() {
  return (
    <input
      type="number"
      name="quantity"
      min="0"
      onWheel={(e) => e.currentTarget.blur()}
    />
  );
}

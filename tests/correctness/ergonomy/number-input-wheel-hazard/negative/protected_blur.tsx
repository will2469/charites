import React from "react";

export function QuantityInput() {
  return (
    <input
      type="number"
      name="quantity"
      min="1"
      onWheel={(e) => e.currentTarget.blur()}
    />
  );
}

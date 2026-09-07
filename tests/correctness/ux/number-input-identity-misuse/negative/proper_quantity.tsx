import React from "react";

export function ProperQuantityInput() {
  return (
    <input
      type="number"
      name="item_count"
      min="1"
      onWheel={(e) => e.currentTarget.blur()}
    />
  );
}

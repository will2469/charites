import React from "react";

export function UnboundedQuantityForm() {
  return (
    <div>
      <label htmlFor="quantity">Quantity</label>
      <input
        type="number"
        id="quantity"
        name="quantity"
        onWheel={(e) => e.currentTarget.blur()}
      />
    </div>
  );
}

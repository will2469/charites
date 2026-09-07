import React from "react";

export function QualifiedQuantitySafeguards() {
  return (
    <div>
      <input
        type="number"
        name="order_quantity"
        min="1"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="account_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="phone_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="postal_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="invoice_total"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="batch_size"
        min="1"
        onWheel={(e) => e.currentTarget.blur()}
      />
    </div>
  );
}

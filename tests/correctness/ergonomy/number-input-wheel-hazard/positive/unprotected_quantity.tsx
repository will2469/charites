import React from "react";

export function OrderForm() {
  return (
    <form>
      <label htmlFor="quantity">Quantity</label>
      <input type="number" id="quantity" name="quantity" />
    </form>
  );
}

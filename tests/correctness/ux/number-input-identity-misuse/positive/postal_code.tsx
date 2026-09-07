import React from "react";

export function AddressForm() {
  return (
    <div>
      <label htmlFor="postal_code">Postal Code</label>
      <input type="number" id="postal_code" name="postal_code" />
    </div>
  );
}

import React from "react";

export function ProperPostalInput() {
  return (
    <input
      type="text"
      inputMode="numeric"
      pattern="[0-9]*"
      name="postal_code"
      autoComplete="postal-code"
    />
  );
}

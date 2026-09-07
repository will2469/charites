import React from "react";

export function IdentityFieldsYield() {
  return (
    <div>
      {/* Identity fields should NOT produce missing-bounds diagnostics because they yield to identity-misuse */}
      <input type="number" name="postal_code" />
      <input type="number" name="order_id" />
      <input type="number" name="account_number" />
      <input type="number" name="nik" />
    </div>
  );
}

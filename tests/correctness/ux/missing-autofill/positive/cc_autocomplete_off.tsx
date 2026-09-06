import React from 'react';

export function PaymentFormWithDisabledAutofill() {
  return (
    <form className="space-y-4">
      <input
        type="text"
        name="cc-number"
        autoComplete="off"
        placeholder="Nomor Kartu Kredit"
        className="border rounded px-3 py-2"
      />
    </form>
  );
}

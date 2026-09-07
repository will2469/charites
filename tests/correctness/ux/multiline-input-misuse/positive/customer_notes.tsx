import React from 'react';

export function OrderNotesForm() {
  return (
    <div className="form-group">
      <label htmlFor="customer-notes">Catatan Pembeli</label>
      <input id="customer-notes" name="customer_notes" />
    </div>
  );
}

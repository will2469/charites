import React from 'react';

export function ShortDescriptionForm() {
  return (
    <div>
      <label htmlFor="desc">Deskripsi Singkat</label>
      <input id="desc" name="description" placeholder="Short description" />
    </div>
  );
}

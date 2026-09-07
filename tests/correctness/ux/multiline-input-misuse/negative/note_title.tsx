import React from 'react';

export function NoteTitleForm() {
  return (
    <div>
      <label htmlFor="note-title">Judul Catatan</label>
      <input id="note-title" name="note_title" placeholder="Tuliskan judul catatan..." />
    </div>
  );
}

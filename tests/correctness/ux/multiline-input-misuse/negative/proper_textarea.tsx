import React from 'react';

export function CompliantNoteForm() {
  return (
    <div className="space-y-2">
      <label htmlFor="keterangan">Keterangan</label>
      <Textarea id="keterangan" name="keterangan" rows={3} placeholder="Tuliskan keterangan lengkap..." />
    </div>
  );
}

function Textarea(props: React.TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return <textarea {...props} />;
}

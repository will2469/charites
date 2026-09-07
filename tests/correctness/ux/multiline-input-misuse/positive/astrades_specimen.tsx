import React from 'react';

export function AstradesSpecimen() {
  return (
    <div className="space-y-2">
      <label htmlFor="keterangan-input">Keterangan (Opsional)</label>
      <Input id="keterangan-input" placeholder="Catatan tambahan (bila ada)" />
    </div>
  );
}

function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return <input {...props} />;
}

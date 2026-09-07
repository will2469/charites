import React from 'react';

export function AstradesDialog() {
  return (
    <div className="p-4 bg-card rounded-xl">
      <p className="font-semibold text-foreground">
        Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?
      </p>
      <br />
      <br />
      <p className="text-sm text-destructive">
        <strong>PERINGATAN:</strong> Tindakan ini bersifat permanen dan tidak dapat dibatalkan setelah konfirmasi.
      </p>
    </div>
  );
}

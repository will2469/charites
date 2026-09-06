import React from "react";

export function TrapModal() {
  return (
    <div role="dialog" aria-modal="true" aria-labelledby="modal-title" className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-background p-6 rounded-xl">
        <h2 id="modal-title" className="text-lg font-semibold">Konfirmasi Hapus</h2>
        <p className="text-sm text-muted-foreground mt-2">Apakah Anda yakin ingin menghapus data ini?</p>
      </div>
    </div>
  );
}

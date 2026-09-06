import React from "react";

export function CompliantModal() {
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Escape") {
      console.log("close");
    }
  };

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      onKeyDown={handleKeyDown}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    >
      <div className="bg-background p-6 rounded-xl relative">
        <button
          type="button"
          onClick={() => console.log("close")}
          aria-label="Tutup modal"
          className="size-11 absolute top-2 right-2 flex items-center justify-center"
        >

        </button>
        <h2 id="modal-title" className="text-lg font-semibold">Konfirmasi Hapus</h2>
        <p className="text-sm text-muted-foreground mt-2">Apakah Anda yakin ingin menghapus data ini?</p>
      </div>
    </div>
  );
}

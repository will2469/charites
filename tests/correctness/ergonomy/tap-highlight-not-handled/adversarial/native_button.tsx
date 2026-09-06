import React from "react";

export function StandardButton({ onConfirm }: { onConfirm: () => void }) {
  // Native buttons are naturally handled by the browser engine and should NOT be flagged
  return (
    <button
      type="button"
      onClick={onConfirm}
      className="h-11 px-4 bg-primary text-primary-foreground rounded-xl font-medium"
    >
      Konfirmasi Permohonan
    </button>
  );
}

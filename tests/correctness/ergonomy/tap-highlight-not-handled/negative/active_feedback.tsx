import React from "react";

export function ActiveFeedbackCard({ onSelect }: { onSelect: () => void }) {
  return (
    // Compliant: deliberate active tactile feedback with suppressed grey highlight
    <div
      role="button"
      tabIndex={0}
      onClick={onSelect}
      className="p-4 bg-card border rounded-2xl cursor-pointer active:bg-muted/60 active:scale-[0.99] transition-transform [-webkit-tap-highlight-color:transparent]"
    >
      <h3>Layanan Permohonan KTP</h3>
      <p>Klik untuk memilih jenis permohonan identitas</p>
    </div>
  );
}

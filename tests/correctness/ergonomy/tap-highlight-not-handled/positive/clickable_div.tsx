import React from "react";

export function CustomActionCard({ onSelect }: { onSelect: () => void }) {
  return (
    // Non-native clickable div without active feedback or tap-highlight styling
    <div
      role="button"
      tabIndex={0}
      onClick={onSelect}
      className="p-4 bg-card border rounded-2xl cursor-pointer"
    >
      <h3>Layanan Permohonan KTP</h3>
      <p>Klik untuk memilih jenis permohonan identitas</p>
    </div>
  );
}

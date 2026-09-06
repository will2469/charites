import React from "react";

export function ResponsiveForm() {
  return (
    // Compliant: min-h-dvh adapts smoothly to keyboard resizing
    <div className="min-h-dvh flex flex-col justify-between pb-safe">
      <input type="text" placeholder="Nama Lengkap" />
      <button className="sticky bottom-4 w-full py-3 bg-primary text-white">Simpan</button>
    </div>
  );
}

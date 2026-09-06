import React from "react";

export function RigidForm() {
  return (
    // Violation: h-screen rigid viewport locks container while housing text input and fixed bottom button
    <div className="fixed inset-0 h-screen flex flex-col justify-between">
      <input type="text" placeholder="Nama Lengkap" />
      <button className="fixed bottom-0 w-full py-3 bg-primary text-white">Simpan</button>
    </div>
  );
}

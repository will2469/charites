import React from "react";

export function HardwareNotes() {
  // Catatan: jangan gunakan navigator.deviceMemory secara langsung
  // karena tidak didukung oleh Apple Safari dan Mozilla Firefox
  return (
    <div className="p-4">
      <p>Aplikasi dioptimalkan untuk perangkat ramah memori.</p>
    </div>
  );
}

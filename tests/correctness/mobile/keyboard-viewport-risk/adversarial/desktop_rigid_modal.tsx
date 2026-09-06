import React from "react";

export function DesktopRigidModal() {
  return (
    // Adversarial: Desktop only (hidden md:flex) is immune from mobile keyboard hazards
    <div className="hidden md:flex h-screen flex-col justify-between">
      <input type="text" placeholder="Admin Query" />
      <button className="fixed bottom-0">Eksekusi</button>
    </div>
  );
}

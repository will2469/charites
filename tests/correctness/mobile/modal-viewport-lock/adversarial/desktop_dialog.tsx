import React from "react";

export function DesktopDialog() {
  return (
    // Adversarial: Desktop only modal (hidden md:flex) is immune
    <div role="dialog" className="hidden md:flex fixed inset-0 overflow-hidden items-center justify-center">
      <div className="bg-card p-6">
        <h2>Desktop Admin Modal</h2>
      </div>
    </div>
  );
}

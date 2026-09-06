import React from "react";

export function DesktopOnlyBar() {
  return (
    // Adversarial: Desktop only bottom bar (hidden md:flex) is immune
    <div className="min-h-screen bg-background">
      <main className="p-4">
        <p>Dashboard desktop...</p>
      </main>
      <div className="hidden md:flex fixed bottom-0 inset-x-0 h-16 bg-card border-t">
        <span>Desktop Tools</span>
      </div>
    </div>
  );
}

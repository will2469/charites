import React from "react";

export function DisabledButton() {
  return (
    // Adversarial: Explicitly disabled button under pointer-events-none container should NOT trigger false-positive
    <div className="pointer-events-none opacity-50 p-4">
      <button disabled className="bg-muted text-muted-foreground px-4 py-2">
        Aksi Terkunci
      </button>
    </div>
  );
}

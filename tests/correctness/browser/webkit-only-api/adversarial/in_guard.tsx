import React from "react";

export function GuardedFullscreen({ target }: { target: HTMLElement | null }) {
  const triggerFullscreen = () => {
    // Compliant: protected by 'in' check
    if (target && "webkitRequestFullscreen" in target) {
      (target as any).webkitRequestFullscreen();
    }
  };

  return (
    <button type="button" onClick={triggerFullscreen}>
      Fullscreen Guarded
    </button>
  );
}

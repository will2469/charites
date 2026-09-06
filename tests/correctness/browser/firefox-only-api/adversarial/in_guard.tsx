import React from "react";

export function GuardedMozElement({ target }: { target: HTMLElement | null }) {
  const triggerFullscreen = () => {
    // Compliant: protected by 'in' check
    if (target && "mozRequestFullScreen" in target) {
      (target as any).mozRequestFullScreen();
    }
  };

  return (
    <button type="button" onClick={triggerFullscreen}>
      Fullscreen Guarded Moz
    </button>
  );
}

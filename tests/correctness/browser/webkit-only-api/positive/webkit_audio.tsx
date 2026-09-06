import React from "react";

export function AudioPlayer() {
  const playBeep = () => {
    // Direct invocation without W3C AudioContext fallback
    const ctx = new webkitAudioContext();
    const osc = ctx.createOscillator();
    osc.start();
  };

  return (
    <button type="button" onClick={playBeep} className="h-11 px-4 bg-primary text-primary-foreground">
      Putar Nada
    </button>
  );
}

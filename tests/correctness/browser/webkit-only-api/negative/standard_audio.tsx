import React from "react";

export function StandardAudioPlayer() {
  const initAudio = () => {
    // Compliant: standard AudioContext with webkit fallback chain
    const AudioCtx = window.AudioContext || window.webkitAudioContext;
    if (AudioCtx) {
      const ctx = new AudioCtx();
      ctx.resume();
    }
  };

  return (
    <button type="button" onClick={initAudio} className="h-11 px-4 bg-primary">
      Inisialisasi Audio
    </button>
  );
}

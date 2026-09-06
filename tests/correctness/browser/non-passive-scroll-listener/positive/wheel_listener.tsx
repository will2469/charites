import React, { useEffect } from "react";

// P2: Wheel listener attached in useEffect without passive: true
export function ZoomableCanvas() {
  useEffect(() => {
    const handleWheel = (e: WheelEvent) => {
      console.log("Delta:", e.deltaY);
    };

    window.addEventListener("wheel", handleWheel);
    return () => window.removeEventListener("wheel", handleWheel);
  }, []);

  return <div className="canvas-container">Canvas Content</div>;
}

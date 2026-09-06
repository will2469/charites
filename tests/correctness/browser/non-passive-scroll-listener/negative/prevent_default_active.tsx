import React, { useEffect } from "react";

// N2: Listener calls e.preventDefault() to cancel scrolling behavior
export function ModalLock() {
  useEffect(() => {
    const blockTouch = (e: TouchEvent) => {
      e.preventDefault();
    };

    window.addEventListener("touchmove", blockTouch);
    return () => window.removeEventListener("touchmove", blockTouch);
  }, []);

  return <div className="modal-backdrop">Modal Active</div>;
}

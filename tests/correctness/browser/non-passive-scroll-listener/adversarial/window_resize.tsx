import React, { useEffect } from "react";

// A2: Window resize or orientationchange events do not require passive: true
export function ResponsiveWatcher() {
  useEffect(() => {
    const onResize = () => {
      console.log("Window resized:", window.innerWidth);
    };

    window.addEventListener("resize", onResize);
    return () => window.removeEventListener("resize", onResize);
  }, []);

  return <div>Responsive Watcher Active</div>;
}

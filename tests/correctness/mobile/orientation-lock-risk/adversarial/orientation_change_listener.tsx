import React, { useEffect } from "react";

export function OrientationListener() {
  useEffect(() => {
    // Adversarial: Listening to orientation change is valid telemetry/layout sync and should NOT trigger lock warning
    const handleChange = () => {
      console.log("Orientation changed to", screen.orientation.type);
    };
    if (screen.orientation) {
      screen.orientation.addEventListener("change", handleChange);
    }
    return () => {
      if (screen.orientation) {
        screen.orientation.removeEventListener("change", handleChange);
      }
    };
  }, []);

  return <div>Monitoring Orientasi</div>;
}

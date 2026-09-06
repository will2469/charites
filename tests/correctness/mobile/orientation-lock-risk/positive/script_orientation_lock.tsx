import React, { useEffect } from "react";

export function LockComponent() {
  useEffect(() => {
    // Violation: programmatic orientation lock restricts mobile accessibility
    if (screen.orientation && screen.orientation.lock) {
      screen.orientation.lock("portrait").catch(() => {});
    }
  }, []);

  return <div>Aplikasi Desa</div>;
}

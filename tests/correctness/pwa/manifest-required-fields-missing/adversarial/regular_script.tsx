import React from "react";

export function RegularScript() {
  return (
    <div>
      {/* Non-manifest script should never trigger PWA manifest rules */}
      <script type="text/javascript" src="/analytics.js"></script>
    </div>
  );
}

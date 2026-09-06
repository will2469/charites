import React from "react";

export function IncompleteManifest() {
  return (
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital"
      })}
    </script>
  );
}

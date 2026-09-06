import React from "react";

export function ManifestWithoutIcons() {
  return (
    // If icons are missing entirely, required-fields-missing reports it.
    // icon-maskable-missing must not double-flag.
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital",
        start_url: "/",
        display: "standalone"
      })}
    </script>
  );
}

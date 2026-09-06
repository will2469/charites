import React from "react";

export function MissingStartURL() {
  return (
    // Missing start_url is flagged by manifest-required-fields-missing,
    // start-url-inconsistency should not flag when empty.
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital",
        display: "standalone",
        icons: [{ src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }]
      })}
    </script>
  );
}

import React from "react";

export function CompleteManifest() {
  return (
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital",
        short_name: "Desa",
        start_url: "/",
        display: "standalone",
        icons: [
          { src: "/icon-192.png", sizes: "192x192", type: "image/png" },
          { src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
        ]
      })}
    </script>
  );
}

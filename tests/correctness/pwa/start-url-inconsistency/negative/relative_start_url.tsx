import React from "react";

export function RelativeStartURL() {
  return (
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital",
        start_url: "/",
        display: "standalone",
        icons: [
          { src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
        ]
      })}
    </script>
  );
}

import React from "react";

export function InsecureStartURL() {
  return (
    <script type="application/manifest+json">
      {JSON.stringify({
        name: "Desa Digital",
        start_url: "http://desa.org/app",
        display: "standalone",
        icons: [
          { src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
        ]
      })}
    </script>
  );
}

import React from "react";

export function PwaHeadWithoutAppleMeta() {
  return (
    <head>
      <title>Layanan Desa</title>
      <link rel="manifest" href="/manifest.webmanifest" />
      {/* Missing apple-mobile-web-app-capable and apple-touch-icon */}
    </head>
  );
}

import React from "react";

export function CompleteApplePwaHead() {
  return (
    <head>
      <title>Layanan Desa</title>
      <link rel="manifest" href="/manifest.webmanifest" />
      <meta name="apple-mobile-web-app-capable" content="yes" />
      <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
    </head>
  );
}

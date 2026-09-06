import React from "react";

export function StandardWebHead() {
  return (
    // Standard web document head without manifest (non-PWA)
    // Should NOT trigger Apple PWA standalone meta warnings.
    <head>
      <title>Blog Pribadi Warga</title>
      <meta name="description" content="Artikel seputar desa digital" />
    </head>
  );
}

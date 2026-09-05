import React from "react";

// N2: Responsive <picture> element with dark mode source
export function ResponsiveLogo() {
  return (
    <picture>
      <source media="(prefers-color-scheme: dark)" srcSet="/images/logo-dark.svg" />
      <img src="/images/logo-black.svg" alt="Company Logo" />
    </picture>
  );
}

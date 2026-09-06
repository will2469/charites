import React from "react";

export function AdversarialErrorInputs() {
  return (
    <div>
      {/* Normal inputs without aria-invalid */}
      <input id="search" type="search" placeholder="Cari..." />
      {/* Non-input div with aria-invalid */}
      <div aria-invalid="true">Container with invalid status</div>
    </div>
  );
}

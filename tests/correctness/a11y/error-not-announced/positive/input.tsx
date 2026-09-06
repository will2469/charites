import React from "react";

export function UnannouncedErrorInputs() {
  return (
    <div>
      <input
        id="email"
        aria-invalid="true"
        className="border-destructive"
        placeholder="nama@domain.com"
      />
    </div>
  );
}

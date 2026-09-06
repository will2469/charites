import React from "react";

export function CompliantErrorInputs() {
  return (
    <div>
      <input
        id="email"
        aria-invalid="true"
        aria-describedby="email-error"
        className="border-destructive"
        placeholder="nama@domain.com"
      />
      <span id="email-error" role="alert">Email format is invalid</span>
      <input
        id="name"
        aria-invalid="false"
        placeholder="Full name"
      />
    </div>
  );
}

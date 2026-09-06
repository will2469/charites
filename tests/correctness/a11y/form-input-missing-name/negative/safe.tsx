import React from "react";

export function CompliantNamedInputs() {
  return (
    <div>
      <input
        id="user-email"
        name="email"
        type="email"
        className="w-full border px-3 py-2"
        placeholder="nama@domain.com"
      />
      <textarea
        id="user-bio"
        name="bio"
        className="border p-2"
      />
    </div>
  );
}

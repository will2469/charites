import React from "react";

export function CompliantLabelInputs() {
  return (
    <div>
      <label htmlFor="user-email">Email</label>
      <input
        id="user-email"
        type="email"
        placeholder="nama@domain.com"
        className="border px-3 py-2"
      />
      <input
        id="search-query"
        name="query"
        type="search"
        placeholder="Cari..."
        aria-label="Pencarian"
        className="border px-3 py-2"
      />
    </div>
  );
}

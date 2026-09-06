import React from "react";

export function MismatchedLabels() {
  return (
    <div>
      <label htmlFor="user_id">User ID</label>
      <input id="userId" className="border px-3 py-2" />
    </div>
  );
}

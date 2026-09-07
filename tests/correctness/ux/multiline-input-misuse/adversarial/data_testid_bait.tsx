import React from "react";

export function DataTestIdBait() {
  return (
    <div>
      {/* data-testid is deliberately excluded from v1 semantic detection */}
      <input
        data-testid="customer-notes-textarea-input"
        name="username"
        placeholder="Enter your username"
      />
    </div>
  );
}

import React from "react";

export function PasswordNotes() {
  return (
    <div>
      {/* Password inputs should not trigger multiline misuse even if named 'notes' */}
      <input type="password" name="notes" placeholder="Enter secure notes" />
    </div>
  );
}

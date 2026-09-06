import React from 'react';

export function NonPIITokenInput() {
  return (
    <form className="space-y-4">
      <input
        type="text"
        name="auth_token_hex"
        placeholder="Auth Token (read-only hash)"
        className="border rounded px-3 py-2 font-mono"
      />
    </form>
  );
}

import React from "react";

// A2: Form with screen-reader-only labels and hidden inputs without hover interactions
export function AccessibleForm() {
  return (
    <form className="space-y-4">
      <label htmlFor="token" className="sr-only">
        Security Token
      </label>
      <input type="hidden" name="token" value="abc-123" />
      <button type="submit" className="px-4 py-2 bg-slate-800 text-white rounded">
        Submit
      </button>
    </form>
  );
}

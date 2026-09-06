import React from "react";

export function StrippedOutlineButtons() {
  return (
    <div>
      <button className="outline-none bg-primary text-white px-4 py-2">
        Submit
      </button>
      <a href="/login" className="focus:outline-none text-blue-500">
        Login
      </a>
    </div>
  );
}

import React from "react";

export function CompliantFocusButtons() {
  return (
    <div>
      <button className="outline-none focus-visible:ring-2 focus-visible:ring-primary bg-primary text-white">
        Save
      </button>
      <a href="/login" className="focus:outline-none focus-visible:ring-2 text-blue-500">
        Login
      </a>
      <button className="focus:outline-none focus:ring-2">
        Profile
      </button>
    </div>
  );
}

import React from "react";

export function SemanticFocusComponents() {
  return (
    <div>
      <button className="focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">Save</button>
      <input className="ring-ring outline-none" />
    </div>
  );
}

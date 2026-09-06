import React from "react";

// P2: Delete icon hidden by default and shown only on group-hover without keyboard/touch focus
export function CardAction() {
  return (
    <div className="group flex items-center justify-between p-3 border">
      <span>File item</span>
      <a
        href="/delete"
        className="hidden group-hover:block text-red-600 font-semibold"
      >
        Delete
      </a>
    </div>
  );
}

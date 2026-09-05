import React from "react";

// P2: Parent opacity-75 with child text-white/60
export function Notification() {
  return (
    <section className="opacity-75 p-4">
      <span className="text-white/60 font-medium">Alert text</span>
    </section>
  );
}

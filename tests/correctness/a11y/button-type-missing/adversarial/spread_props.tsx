import React from "react";

export function AdversarialButtons(props: Record<string, unknown>) {
  return (
    <form className="space-y-4">
      {/* Spread props forwarding */}
      <button {...props} className="h-11 px-4 text-sm">
        Forwarded
      </button>

      {/* Decorative / hidden button */}
      <button aria-hidden="true" className="h-11 px-4 text-sm">
        Hidden
      </button>
    </form>
  );
}

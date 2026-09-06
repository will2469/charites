import React from "react";

export function AdversarialTouchTargets() {
  return (
    <div>
      {/* Non-interactive spans and divs with small sizes should not trigger */}
      <span className="h-4 w-4">icon</span>
      <div className="size-6">badge</div>
      <p>
        Visit <a href="/privacy" className="text-sm underline">privacy page</a> in text.
      </p>
    </div>
  );
}

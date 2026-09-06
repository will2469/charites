import React from "react";

export function AdversarialAriaDialog(props: Record<string, unknown>, isOpen: boolean) {
  return (
    <div>
      {/* Spread props */}
      <div role="dialog" {...props} className="fixed inset-0">
        <p>Dialog with props</p>
      </div>

      {/* Dynamic aria-modal expression */}
      <div role="dialog" aria-modal={isOpen} aria-label="Dynamic Modal" className="fixed inset-0">
        <p>Dynamic modal</p>
      </div>

      {/* Non-dialog role */}
      <div role="region" aria-label="Main region">
        <p>Region content</p>
      </div>
    </div>
  );
}

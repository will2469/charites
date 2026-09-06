import React from "react";

export function AdversarialModal(props: Record<string, unknown>) {
  return (
    <div>
      {/* Spread props forwarding */}
      <div role="dialog" {...props} className="fixed inset-0">
        <p>Dynamic modal</p>
      </div>

      {/* Hidden overlay */}
      <div role="dialog" aria-hidden="true" className="hidden">
        <p>Decorative hidden dialog</p>
      </div>
    </div>
  );
}

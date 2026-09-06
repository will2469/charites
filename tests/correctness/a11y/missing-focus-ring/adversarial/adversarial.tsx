import React from "react";

export function AdversarialFocusElements() {
  return (
    <div>
      {/* Non-interactive div with outline-none is safe */}
      <div className="outline-none p-4">
        Container content
      </div>
      {/* Interactive element using replacement outline */}
      <button className="outline-none focus-visible:outline-2 focus-visible:outline-blue-500">
        Custom outline
      </button>
    </div>
  );
}

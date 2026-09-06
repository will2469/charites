import React from "react";

export function CompliantTouchTargets() {
  return (
    <div>
      <button className="h-11 w-11">Safe Button</button>
      <button className="h-8 w-8 min-h-11 min-w-11">Compensated Button</button>
      <p>
        Please check <a href="/terms" className="text-blue-600">terms and conditions</a>.
      </p>
      <div className="h-6 w-6">Non-interactive div</div>
    </div>
  );
}

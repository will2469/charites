import React from "react";

export function AdversarialLabels() {
  const dynamicId = "dynamic-field";
  return (
    <div>
      {/* Implicit label wrapping input without htmlFor */}
      <label>
        Username <input id="nav-username" name="username" type="text" />
      </label>
      {/* Dynamic htmlFor expression */}
      <label htmlFor={dynamicId}>Field</label>
      <input id={dynamicId} />
    </div>
  );
}

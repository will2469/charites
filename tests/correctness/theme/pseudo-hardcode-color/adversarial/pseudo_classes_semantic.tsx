import React from "react";

export function PseudoElementsSemantic() {
  return (
    <div className="before:bg-accent after:bg-primary before:text-foreground">
      <span>Semantic Pseudo Elements</span>
    </div>
  );
}

import React from "react";

export function SafeComponent() {
  return (
    <div className="p-4 bg-card text-card-foreground border-border">
      <span className="bg-[var(--primary)] text-primary-foreground">Safe</span>
    </div>
  );
}

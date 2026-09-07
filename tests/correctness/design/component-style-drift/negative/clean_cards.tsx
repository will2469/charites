import React from "react";

// N2: Compliant card container with clean semantic tokens and outer radius
export function CleanCards() {
  return (
    <Card className="rounded-xl bg-card text-card-foreground p-6 shadow-sm border border-border">
      <h2 className="text-xl font-semibold">Card Title</h2>
      <p className="text-sm text-muted-foreground">Card content with compliant tokens.</p>
    </Card>
  );
}

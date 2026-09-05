import React from "react";

export function UserProfile() {
  return (
    <div className="flex items-center justify-between p-4 rounded-lg bg-background text-foreground">
      <div className="space-y-1">
        <h4 className="text-sm font-medium leading-none text-foreground">User Name</h4>
        <p className="text-sm text-muted">user@example.com</p>
      </div>
      <button className="inline-flex items-center justify-center rounded-md text-sm font-medium border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground">
        Edit
      </button>
    </div>
  );
}

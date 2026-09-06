import React from "react";

export function CompliantSpacing() {
  return (
    <div>
      <div className="flex gap-2">
        <button className="h-11 w-11">Edit</button>
        <button className="h-11 w-11">Delete</button>
      </div>
      <div className="flex gap-4">
        <button className="h-11 w-11">Cancel</button>
        <button className="h-11 w-11">Save</button>
      </div>
      <div className="flex">
        <button className="h-11 w-11">Only One Button</button>
        <span>Informative text</span>
      </div>
    </div>
  );
}

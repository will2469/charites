import React from "react";

export function DynamicTypeNotes({ customType }: { customType: string }) {
  return (
    <div>
      {/* Dynamic type inputs should be treated as Unknown type, not triggering text rule */}
      <input type={customType} name="notes" placeholder="Enter notes" />
    </div>
  );
}

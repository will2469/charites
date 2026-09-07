import React from "react";

export function NumberNotesCount() {
  return (
    <div>
      {/* Number inputs with count quantifier should not trigger multiline misuse */}
      <input
        type="number"
        name="notes_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
    </div>
  );
}

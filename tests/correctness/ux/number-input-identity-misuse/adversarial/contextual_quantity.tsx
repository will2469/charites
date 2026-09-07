import React from "react";

export function ContextualQuantityInput() {
  return (
    <div>
      <input
        type="number"
        name="total_rt_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
      <input
        type="number"
        name="wa_count"
        min="0"
        onWheel={(e) => e.currentTarget.blur()}
      />
    </div>
  );
}

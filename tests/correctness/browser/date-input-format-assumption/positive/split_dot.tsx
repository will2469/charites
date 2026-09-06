import React from "react";

// P2: Splitting native date input by dot in onChange
export function DatePicker() {
  return (
    <input
      type="date"
      onChange={(e) => {
        const parts = e.target.value.split('.');
        console.log(parts);
      }}
    />
  );
}

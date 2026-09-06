import React from "react";

// N2: Using valueAsDate or standard Date constructor
export function CompliantDateInput() {
  return (
    <input
      type="date"
      onChange={(e) => {
        const dateObj = (e.target as HTMLInputElement).valueAsDate;
        console.log(dateObj);
      }}
    />
  );
}

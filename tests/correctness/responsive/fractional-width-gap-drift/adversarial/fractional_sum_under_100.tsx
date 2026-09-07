import React from "react";

export function FractionalSumUnder100() {
  return (
    // sum is 50% < 100%, fits easily with gap even under wrap or shrink-0
    <div className="flex flex-wrap gap-4">
      <div className="w-1/4 shrink-0">A</div>
      <div className="w-1/4">B</div>
    </div>
  );
}

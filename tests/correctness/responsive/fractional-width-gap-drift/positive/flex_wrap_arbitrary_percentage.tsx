import React from "react";

export function FlexWrapArbitraryPercentage() {
  return (
    <div className="flex flex-wrap gap-4">
      <div className="w-[60%]">Kolom Kiri</div>
      <div className="w-[40%]">Kolom Kanan</div>
    </div>
  );
}

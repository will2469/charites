import React from "react";

export function CustomBrandedBadge() {
  return (
    <div className="flex gap-2">
      <span className="p-2 bg-[#123456]/10 text-[#ff0000]/20 border border-[#abcdef]/5">
        Arbitrary Hex Colors (Handled by a separate rule, not theme.hardcode-opacity-color)
      </span>
    </div>
  );
}

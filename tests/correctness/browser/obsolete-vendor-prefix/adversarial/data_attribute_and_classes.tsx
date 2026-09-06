import React from "react";

// A1: Data attributes, class names, and text content containing prefix substrings
export function PrefixLikeMarkup() {
  return (
    <div
      data-testid="-webkit-transform"
      data-vendor-prefix="-moz-box-sizing"
      className="line-clamp-2 prefix-moz box-border"
    >
      <span>Legacy properties like -webkit-border-radius should not trigger in text</span>
    </div>
  );
}

import React from "react";

export function MobileDetectBanner() {
  const isMobile = typeof navigator !== "undefined" && /android|iphone|ipad/i.test(navigator.userAgent);

  return (
    <div className="p-4 bg-muted text-sm">
      {isMobile ? <span>Mobile view active</span> : <span>Desktop view active</span>}
    </div>
  );
}

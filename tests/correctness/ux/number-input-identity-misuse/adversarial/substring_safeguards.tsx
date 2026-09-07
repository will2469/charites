import React from "react";

export function SubstringSafeguards() {
  return (
    <div>
      <input type="text" name="start_date" />
      <input type="number" name="alert_count" min="0" onWheel={(e) => e.currentTarget.blur()} />
      <input type="text" name="php_version" />
      <input type="text" name="short_code" />
    </div>
  );
}

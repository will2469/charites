import React from "react";

export function AbbreviatedDesc() {
  return (
    <div>
      {/* 'desc' is not a strong multiline token, preventing collision with sorting or internal abbreviations */}
      <input name="sort_desc" defaultValue="descending" />
      <input name="api_desc" />
    </div>
  );
}

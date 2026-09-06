import React from "react";

export function CitizenProfileName() {
  return (
    <div className="space-y-2">
      <label htmlFor="name_input">Nama Lengkap Sesuai KTP</label>
      {/* Standard text input without numeric/tel keywords should NOT be flagged */}
      <input
        id="name_input"
        name="full_name"
        placeholder="Budi Santoso"
        className="h-11 px-3.5 py-2.5 bg-background border border-input rounded-xl text-base"
      />
    </div>
  );
}

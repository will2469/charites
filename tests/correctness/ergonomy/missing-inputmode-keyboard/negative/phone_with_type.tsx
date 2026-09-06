import React from "react";

export function StandardPhoneForm() {
  return (
    <form className="space-y-4">
      <label htmlFor="hp_field" className="text-sm font-medium">Nomor Handphone</label>
      {/* Compliant: has type="tel" and inputMode="tel" */}
      <input
        id="hp_field"
        name="nomor_hp"
        type="tel"
        inputMode="tel"
        autoComplete="tel"
        placeholder="081234567890"
        className="h-11 px-3.5 py-2.5 bg-background border border-input rounded-xl text-base"
      />
    </form>
  );
}

import React from "react";

export function CitizenRegistrationForm() {
  return (
    <form className="space-y-4">
      <label htmlFor="hp_field" className="text-sm font-medium">Nomor Handphone</label>
      {/* Missing type="tel" or inputmode */}
      <input
        id="hp_field"
        name="nomor_hp"
        placeholder="081234567890"
        className="h-11 px-3.5 py-2.5 bg-background border border-input rounded-xl text-base"
      />
    </form>
  );
}

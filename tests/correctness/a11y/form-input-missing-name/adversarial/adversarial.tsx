import React from "react";

export function AdversarialNamedInputs() {
  const register = (name: string) => ({ name });
  return (
    <div>
      {/* Submit/button inputs don't require name or id */}
      <input type="submit" value="Kirim" />
      <input type="button" value="Batal" />
      <input type="reset" value="Reset" />
      {/* react-hook-form spread register */}
      <input {...register("email")} />
    </div>
  );
}

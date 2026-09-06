import React from "react";

export function CompliantButtons() {
  return (
    <div>
      {/* Tombol di luar form sah tanpa type */}
      <button onClick={() => console.log("standalone")} className="h-11 px-4 text-sm font-medium">
        Di Luar Form
      </button>

      {/* Tombol di dalam form dengan type eksplisit */}
      <form className="space-y-4">
        <button type="button" onClick={() => console.log("cancel")} className="h-11 px-4 text-sm font-medium">
          Batal
        </button>
        <button type="reset" className="h-11 px-4 text-sm font-medium">
          Reset
        </button>
        <button type="submit" className="h-11 px-4 text-sm font-medium bg-primary text-primary-foreground">
          Kirim
        </button>
      </form>
    </div>
  );
}

import React from "react";

export function FormWithUntypedButtons() {
  const handleCancel = () => console.log("cancel");
  const handlePreview = () => console.log("preview");

  return (
    <form onSubmit={(e) => e.preventDefault()} className="space-y-4">
      <input id="title" name="title" className="h-11 px-3.5 text-base border rounded-lg" />
      <button onClick={handlePreview} className="h-11 px-4 text-sm font-medium">
        Pratinjau
      </button>
      <button onClick={handleCancel} className="h-11 px-4 text-sm font-medium">
        Batal
      </button>
      <button type="submit" className="h-11 px-4 text-sm font-medium bg-primary text-primary-foreground">
        Kirim
      </button>
    </form>
  );
}

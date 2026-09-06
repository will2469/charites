import React from "react";

// N2: showOpenFilePicker guarded with 'showOpenFilePicker' in window
export function GuardedPicker() {
  const handleOpen = async () => {
    if ("showOpenFilePicker" in window) {
      const files = await (window as any).showOpenFilePicker();
      console.log(files);
    } else {
      document.getElementById("legacy-file-input")?.click();
    }
  };

  return (
    <div>
      <button onClick={handleOpen}>Upload</button>
      <input id="legacy-file-input" type="file" className="hidden" />
    </div>
  );
}

import React from "react";

// P2: Direct invocation of window.showOpenFilePicker without feature guard
export function FilePickerButton() {
  const handleOpen = async () => {
    const handles = await window.showOpenFilePicker();
    console.log(handles);
  };

  return (
    <button onClick={handleOpen} className="px-4 py-2 bg-blue-600 text-white rounded">
      Open File
    </button>
  );
}

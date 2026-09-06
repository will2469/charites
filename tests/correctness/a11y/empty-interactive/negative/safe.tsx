import React from "react";

export function CompliantInteractiveElements() {
  return (
    <div>
      {/* Icon button with explicit aria-label */}
      <button
        type="button"
        onClick={() => console.log("del")}
        aria-label="Hapus dokumen"
        className="size-11 flex items-center justify-center"
      >
        <TrashIcon className="size-5" />
      </button>

      {/* Button with visible text label */}
      <button type="button" className="h-11 px-4 text-sm font-medium">
        <span>Hapus</span>
      </button>

      {/* Link with text and icon */}
      <a href="/settings" className="h-11 px-4 flex items-center gap-2">
        <SettingsIcon className="size-5" />
        <span>Pengaturan</span>
      </a>

      {/* Button with title attribute */}
      <button type="button" title="Unduh laporan" className="size-11 flex items-center justify-center">
        <DownloadIcon className="size-5" />
      </button>
    </div>
  );
}

function TrashIcon(props: { className?: string }) {
  return <svg {...props} />;
}
function SettingsIcon(props: { className?: string }) {
  return <svg {...props} />;
}
function DownloadIcon(props: { className?: string }) {
  return <svg {...props} />;
}

import React from "react";

export function IconButtonViolations() {
  return (
    <div>
      {/* Icon-only button without accessible name */}
      <button onClick={() => console.log("del")} className="size-11 flex items-center justify-center">
        <TrashIcon className="size-5" />
      </button>

      {/* SVG-only link without accessible name */}
      <a href="/settings" className="size-11 flex items-center justify-center">
        <svg className="size-5" viewBox="0 0 24 24">
          <path d="M12 2L2 7l10 5 10-5-10-5z" />
        </svg>
      </a>

      {/* Role button with only an icon */}
      <div role="button" onClick={() => console.log("action")} className="size-11 flex items-center justify-center">
        <EditIcon className="size-5" />
      </div>
    </div>
  );
}

function TrashIcon(props: { className?: string }) {
  return <svg {...props} />;
}

function EditIcon(props: { className?: string }) {
  return <svg {...props} />;
}

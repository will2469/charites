import React, { useRef } from 'react';

export function NativeDialogModal() {
  const dialogRef = useRef<HTMLDialogElement>(null);
  return (
    <dialog ref={dialogRef} className="fixed inset-y-0 right-0 w-96">
      <div className="p-4">
        <h3>Native Dialog Modal</h3>
      </div>
    </dialog>
  );
}

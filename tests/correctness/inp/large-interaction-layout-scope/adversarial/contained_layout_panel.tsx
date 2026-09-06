import React from 'react';

export function ContainedPanel({ isOpen }: { isOpen: boolean }) {
  return (
    <div className={`fixed inset-y-0 right-0 w-96 contain-layout ${isOpen ? "block" : "hidden"}`}>
      <div className="p-4">
        <h3>Contained Layout Panel</h3>
      </div>
    </div>
  );
}

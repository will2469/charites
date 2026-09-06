import React from 'react';

export function UnconstrainedDrawer({ isOpen }: { isOpen: boolean }) {
  return (
    <div className={`fixed inset-y-0 right-0 w-96 ${isOpen ? "block" : "hidden"}`}>
      <div className="p-4">
        <h3>Sidebar Drawer</h3>
      </div>
    </div>
  );
}

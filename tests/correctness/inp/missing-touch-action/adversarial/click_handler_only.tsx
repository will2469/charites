import React from 'react';

export function StandardClickButton({ handleClick }: { handleClick: () => void }) {
  return (
    <button onClick={handleClick} className="px-4 py-2 bg-primary text-white">
      Click Me
    </button>
  );
}

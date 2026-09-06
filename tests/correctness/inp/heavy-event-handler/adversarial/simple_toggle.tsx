import React, { useState } from 'react';

export function SimpleToggle() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <button onClick={() => setIsOpen(!isOpen)}>
      Toggle Menu
    </button>
  );
}
